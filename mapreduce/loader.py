from collections import namedtuple
from csv import DictReader
from datetime import datetime

from cassandra.cluster import Cluster

Record = namedtuple("Record", ["invoice", "stock_code", "date", "description", "unit_price", "customer", "country"])

OrderPart = namedtuple("OrderPart", ["amount", "good_description", "good_price", "good_code"])
GoodStat = namedtuple("GoodStat", ["price", "bought", "total_price"])
if __name__ == '__main__':

    cluster = Cluster(['localhost'], port=9088)
    session = cluster.connect("shop")

    cluster.register_user_type('shop', 'order_part', OrderPart)

    print(session.execute("SELECT release_version FROM system.local").one())

    goods = {}
    orders = {}

    orders_stats_by_keyword = {}
    orders_stats_by_county = {}
    orders_stats_by_good = {}
    with open("./data/data.csv") as f:
        r = DictReader(f)
        for row in r:
            invoice_no = row["InvoiceNo"]
            stock_code = row["StockCode"]
            invoice_date = datetime.strptime(row["InvoiceDate"], "%m/%d/%Y %H:%M")
            description = row["Description"]
            quantity = int(row["Quantity"])
            if quantity < 0:
                continue
            unit_price = float(row["UnitPrice"])
            try:
                customer = int(row["CustomerID"])
            except ValueError:
                continue
            country = row["Country"]
            # rows.append(Record(invoice_no, stock_code, invoice_date, description, unit_price, customer, country))

            session.execute("INSERT INTO orders_raw (customer_id, invoice_no, destination_country, created_at,stock_code, quantity, unit_price) VALUES (%s, %s, %s, %s, %s, %s, %s)", (customer, invoice_no, country, invoice_date, stock_code, quantity, unit_price))
            goods[stock_code] = (unit_price, description)
            if invoice_no not in orders:
                orders[invoice_no] = (invoice_date, customer, country, [])

            o = orders[invoice_no][3]
            o.append((stock_code, unit_price, description, quantity))

            date = invoice_date.date()
            for keyword in description.split():
                keyword = keyword.strip().lower()
                dates = orders_stats_by_keyword.setdefault(keyword, dict())

                amount, total_price, lgoods = dates.setdefault(date, (0, 0, {}))
                amount += 1
                total_price += quantity * unit_price

                if stock_code not in lgoods:
                    lgoods[stock_code] = GoodStat(unit_price, 0, 0)

                _, q, tp = lgoods[stock_code]
                q += quantity
                tp += quantity * unit_price
                lgoods[stock_code] = GoodStat(unit_price, q, tp)

                dates[date] = (amount, total_price, lgoods)

            country_dates = orders_stats_by_county.setdefault(country, dict())
            amount, total_price, lgoods = country_dates.setdefault(date, (0, 0, {}))
            amount += 1
            total_price += quantity * unit_price

            if stock_code not in lgoods:
                lgoods[stock_code] = GoodStat(unit_price, 0, 0)

            _, q, tp = lgoods[stock_code]
            q += quantity
            tp += quantity * unit_price
            lgoods[stock_code] = GoodStat(unit_price, q, tp)

            country_dates[date] = (amount, total_price, lgoods)

            good_dates = orders_stats_by_good.setdefault(stock_code, dict())
            amount, total_price, lgoods = good_dates.setdefault(date, (0, 0, {}))
            amount += 1
            total_price += quantity * unit_price

            if stock_code not in lgoods:
                lgoods[stock_code] = GoodStat(unit_price, 0, 0)

            _, q, tp = lgoods[stock_code]
            q += quantity
            tp += quantity * unit_price
            lgoods[stock_code] = GoodStat(unit_price, q, tp)

            good_dates[date] = (amount, total_price, lgoods)

    # for code, (price, description) in goods.items():
    #     session.execute("INSERT INTO goods (description, price, code) values (%s, %s, %s)", (description, unit_price, code))

    # for invoice, (date, customer, country, o) in orders.items():
    #     session.execute("INSERT INTO orders (customer_id, invoice_no, destination_country, created_at, ordered) VALUES (%s, %s, %s, %s, %s)",
    #                     (customer, invoice, country, date, [OrderPart(quantity, description, unit_price, stock_code) for (stock_code, unit_price, description, quantity) in o]))

    cluster.register_user_type('shop', 'good_stat', GoodStat)
    #
    # print(len(orders_stats_by_keyword))
    # for keyword, sub in orders_stats_by_keyword.items():
    #     for date, (amount, total_price, lgoods) in sub.items():
    #         session.execute("INSERT INTO orders_stats_by_keyword (keyword, date, order_amount, total_price, goods) values (%s, %s, %s, %s, %s)",
    #                         (keyword, date, amount, total_price, lgoods))
    #
    # print(len(orders_stats_by_county))
    # for country, sub in orders_stats_by_county.items():
    #     for date, (amount, total_price, lgoods) in sub.items():
    #         session.execute("INSERT INTO orders_stats_by_county (country, date, order_amount, total_price, goods) values (%s, %s, %s, %s, %s)",
    #                         (country, date, amount, total_price, lgoods))
    #
    # print(len(orders_stats_by_good))
    # for good, sub in orders_stats_by_good.items():
    #     for date, (amount, total_price, lgoods) in sub.items():
    #         session.execute("INSERT INTO orders_stats_by_good (good_code, date, order_amount, total_price, goods) values (%s, %s, %s, %s, %s)",
    #                         (good, date, amount, total_price, lgoods))
