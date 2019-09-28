select count(*)
from invoices;

select count(*)
from invoice_parts;

-- 1. Price of invoice (to show for user)
explain analyse verbose
select sum(unit_price * quantity)
from invoice_parts
where invoice_id = 536385;

-- 2. Month income (for bills or statistics)
explain analyse
select sum(unit_price * quantity)
from invoice_parts ip
         join invoices i on ip.invoice_id = i.id
where date(invoice_date) between '2011-01-01' and '2011-02-01';

-- 3. Goods search (for users)
-- literal match
explain analyse
select id, description
from goods
where description like '%TEA%'
  and price between 10 and 30;

-- 4. similarity match
explain analyse
select id, description
from goods
where to_tsvector(description) @@ to_tsquery('tea')
  and price between 10 and 30;

-- 5. Country chart by invoice count (statistics)
explain analyse
select destination_country, count(*), sum(ip.unit_price * ip.quantity)
from invoices
         join invoice_parts ip on invoices.id = ip.invoice_id
where date(invoice_date) between '2011-01-01' and '2012-01-01'
group by destination_country
order by count(*) desc;

-- 6. Goods sold quantities (to order more or to show recommendations)
explain analyse
select g.id, description, sum(ip.quantity) as sold
from goods g
         join invoice_parts ip on g.id = ip.good_id
         join invoices i on ip.invoice_id = i.id
where date(i.invoice_date) between '2011-01-01' and '2011-04-01'
group by g.id
order by sold desc;

-- 7. "My orders"
explain analyse
select i.id, g.id, ip.quantity, g.description
from invoices i
         join invoice_parts ip on i.id = ip.invoice_id
         join goods g on ip.good_id = g.id
where i.customer_id = 15362;
