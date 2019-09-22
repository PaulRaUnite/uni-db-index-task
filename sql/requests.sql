select count(*)
from invoices;

select count(*)
from invoice_parts;

explain analyse verbose
select sum(unit_price * quantity)
from invoice_parts
where invoice_id = 536385;

explain analyse
select sum(unit_price * quantity)
from invoice_parts ip
         join invoices i on ip.invoice_id = i.id
where date(invoice_date) between '2011-01-01' and '2011-02-01';

explain analyse
select id
from goods
where lower(description) like '%tea%';