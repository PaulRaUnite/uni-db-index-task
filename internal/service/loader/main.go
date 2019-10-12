package loader

import (
	"encoding/csv"
	"io"

	"github.com/shopspring/decimal"

	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/internal/config"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/PaulRaUnite/uni-db-index-task/internal/data"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

type good struct {
	id          int
	latestPrice decimal.Decimal
	updated     bool
}
type impl struct {
	db        *pgdb.DB
	source    *gocsv.Unmarshaller
	customers map[int]struct{}
	goods     map[string]good
	invoices  map[int]struct{}
	countries map[string]int
}

func Run(cfg config.Config, source io.Reader) error {
	unm, err := gocsv.NewUnmarshaller(csv.NewReader(source), data.Record{})
	if err != nil {
		return errors.Wrap(err, "failed to create unmarshaller")
	}
	i := impl{
		db:     cfg.DB(),
		source: unm,
	}
	err = i.preLoad()
	if err != nil {
		return errors.Wrap(err, "failed to preload records")
	}
	err = i.runOnce()
	if err != nil {
		return err
	}
	return errors.Wrap(i.dumpPrices(), "failed to dump prices")
}

func (i *impl) runOnce() error {
	for {
		record, err := i.source.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal record")
		}
		err = i.upsertInvoicePart(record.(data.Record))
		if err != nil {
			return errors.Wrapf(err, "failed to save record in database: %v", record)
		}
	}
}

func (i *impl) preLoad() error {
	var ids []struct{ ID int }
	err := i.db.Select(&ids, squirrel.Select("id").From("customers"))
	if err != nil {
		return err
	}
	i.customers = make(map[int]struct{}, 2*len(ids))
	for _, record := range ids {
		i.customers[record.ID] = struct{}{}
	}

	ids = nil
	err = i.db.Select(&ids, squirrel.Select("id").From("invoices"))
	if err != nil {
		return err
	}
	i.invoices = make(map[int]struct{}, 2*len(ids))
	for _, record := range ids {
		i.invoices[record.ID] = struct{}{}
	}
	var goods []struct {
		ID    int
		Code  string
		Price decimal.Decimal
	}
	err = i.db.Select(&goods, squirrel.Select("id", "code", "price").From("goods"))
	if err != nil {
		return err
	}
	i.goods = make(map[string]good, 2*len(goods))
	for _, record := range goods {
		i.goods[record.Code] = good{id: record.ID, latestPrice: record.Price}
	}
	var countries []struct {
		ID           int    `db:"id"`
		ReadableName string `db:"readable_name"`
	}
	err = i.db.Select(&countries, squirrel.Select("*").From("countries"))
	if err != nil {
		return err
	}
	i.countries = make(map[string]int, 2*len(countries))
	for _, record := range countries {
		i.countries[record.ReadableName] = record.ID
	}
	return nil
}

func (i *impl) dumpPrices() error {
	for _, r := range i.goods {
		if !r.updated {
			continue
		}
		err := i.db.Exec(squirrel.Update("goods").
			Where(squirrel.Eq{"id": r.id}).
			SetMap(map[string]interface{}{"price": r.latestPrice}),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *impl) upsertInvoicePart(record data.Record) error {
	err := i.upsertCustomer(record.CustomerID)
	if err != nil {
		return err
	}
	err = i.upsertInvoice(record)
	if err != nil {
		return err
	}
	goodID, err := i.upsertGood(record.StockCode, record.UnitPrice, record.Description)
	if err != nil {
		return err
	}

	return i.insertInvoicePart(goodID, record)
}

func (i *impl) upsertCountry(country string) (int, error) {
	if id, ok := i.countries[country]; ok {
		return id, nil
	}

	insertQuery := squirrel.Insert("countries").
		SetMap(map[string]interface{}{
			"readable_name": country,
		}).Suffix("RETURNING id")

	var id int
	err := i.db.Get(&id, insertQuery)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert into countries")
	}
	i.countries[country] = id
	return id, nil
}

func (i *impl) insertInvoicePart(goodID int, record data.Record) error {
	query := squirrel.Insert("invoice_parts").
		SetMap(map[string]interface{}{
			"good_id":    goodID,
			"unit_price": record.UnitPrice,
			"quantity":   record.Quantity,
			"invoice_id": record.InvoiceNo,
		})
	err := i.db.Exec(query)
	return errors.Wrap(err, "failed to insert invoice part")
}

const postgrestimestampFormat = "2006-01-02 15:04:05-07"

func (i *impl) upsertInvoice(record data.Record) error {
	if _, ok := i.invoices[record.InvoiceNo]; ok {
		return nil
	}

	countryID, err := i.upsertCountry(record.Country)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert a country: %s", record.Country)
	}
	insertQuery := squirrel.Insert("invoices").
		SetMap(map[string]interface{}{
			"id":                     record.InvoiceNo,
			"customer_id":            record.CustomerID,
			"destination_country_id": countryID,
			"invoice_date":           record.InvoiceDate.Format(postgrestimestampFormat),
		})

	err = i.db.Exec(insertQuery)
	if err != nil {
		return errors.Wrap(err, "failed to insert into invoices")
	}
	i.invoices[record.InvoiceNo] = struct{}{}
	return nil
}

func (i *impl) upsertCustomer(customerID int) error {
	if _, ok := i.customers[customerID]; ok {
		return nil
	}

	query := squirrel.Insert("customers").
		SetMap(map[string]interface{}{
			"id": customerID,
		})

	err := i.db.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to insert customer")
	}
	i.customers[customerID] = struct{}{}
	return nil
}

func (i *impl) upsertGood(code string, price decimal.Decimal, description string) (int, error) {
	g, ok := i.goods[code]
	if ok {
		g.latestPrice = price
		g.updated = true
		i.goods[code] = g
		return g.id, nil
	}

	insertQuery := squirrel.Insert("goods").
		SetMap(map[string]interface{}{
			"code":        code,
			"price":       price,
			"description": description,
		}).Suffix("RETURNING id")

	var id int
	err := i.db.Get(&id, insertQuery)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert into goods")
	}
	i.goods[code] = good{
		id:          id,
		latestPrice: price,
	}
	return id, nil
}
