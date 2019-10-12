create extension pg_trgm;

create index invoice_parts_invoice_id_idx on invoice_parts(invoice_id);
create index invoice_customer_id_idx on invoices(customer_id);
create index invoice_date_date_idx on invoices(date(invoice_date));
create index goods_descr_trigram_idx on goods using gist (description gist_trgm_ops);
create index goods_descr_gin_idx on goods using gin (to_tsvector('english', description));
create index goods_price_idx on goods(price);

drop index invoice_parts_invoice_id_idx;
drop index invoice_customer_id_idx;
drop index invoice_date_date_idx;
drop index goods_descr_trigram_idx;
drop index goods_descr_gin_idx;
drop index goods_price_idx;

CREATE EXTENSION pgstattuple;

SELECT 'invoice_parts_invoice_id_idx', * FROM pgstatindex('invoice_parts_invoice_id_idx')
UNION
SELECT 'invoice_customer_id_idx', * FROM pgstatindex('invoice_customer_id_idx')
UNION
SELECT 'invoice_date_date_idx', * FROM pgstatindex('invoice_date_date_idx')
UNION
SELECT 'goods_price_idx', * FROM pgstatindex('goods_price_idx');
