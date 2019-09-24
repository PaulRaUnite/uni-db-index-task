select count(*)
from invoices;

select count(*)
from invoice_parts;

-- Price of invoice (to show for user)
explain analyse verbose
select sum(unit_price * quantity)
from invoice_parts
where invoice_id = 536385;

-- Month income (for statistics)
explain analyse
select sum(unit_price * quantity)
from invoice_parts ip
         join invoices i on ip.invoice_id = i.id
where date(invoice_date) between '2011-01-01' and '2011-02-01';

-- Goods search (for users)
explain analyse
select id
from goods
where description like '%TEA%' and price between 10 and 30;

explain analyse
select id
from goods
where to_tsvector(description) @@ to_tsquery('TEA');

-- Country chart by invoice count (statistics)
explain analyse
select destination_country, count(*), sum(ip.unit_price*ip.quantity)
from invoices join invoice_parts ip on invoices.id = ip.invoice_id
where date(invoice_date) between '2011-01-01' and '2012-01-01'
group by destination_country order by count(*) desc;

-- Goods chart (statistics)
explain analyse
select g.id, description, sum(ip.quantity*ip.unit_price)/count(i.id) as average_bill
from goods g join invoice_parts ip on g.id = ip.good_id join invoices i on ip.invoice_id = i.id
where date(i.invoice_date) between '2011-01-01' and '2012-01-01'
group by g.id order by average_bill desc;

-- Goods sold quantities
explain analyse
select g.id, description, sum(ip.quantity) as sold
from goods g join invoice_parts ip on g.id = ip.good_id join invoices i on ip.invoice_id = i.id
where date(i.invoice_date) between '2011-01-01' and '2011-04-01'
group by g.id order by sold desc;
