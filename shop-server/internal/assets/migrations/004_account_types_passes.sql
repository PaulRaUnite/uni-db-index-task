-- +migrate Up notransaction
alter table customers
    add column password     varchar(256) not null default 'test',
    add column account_type int default 0,
    add column login        varchar(256);

CREATE UNIQUE INDEX CONCURRENTLY login_unique_idx
    ON customers (login);

alter table customers
    add constraint unique_login unique using index login_unique_idx;

alter table invoice_parts
    drop constraint invoice_parts_invoice_id_fkey,
    add constraint invoice_parts_invoice_id_fkey foreign key (invoice_id) references invoices (id) on delete cascade;

alter table invoices
    drop constraint invoices_customer_id_fkey,
    add constraint invoices_customer_id_fkey foreign key (customer_id) references customers (id) on delete cascade;

delete
from customers
where name = '';

update customers
set login = lower(name);

alter table customers
    alter column login set not null;


-- +migrate Down

alter table customers
    drop column password,
    drop column account_type,
    drop column login;
