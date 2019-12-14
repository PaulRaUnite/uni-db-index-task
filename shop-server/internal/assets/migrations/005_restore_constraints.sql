-- +migrate Up

alter table invoice_parts
    drop constraint invoice_parts_invoice_id_fkey,
    add constraint invoice_parts_invoice_id_fkey foreign key (invoice_id) references invoices (id);

alter table invoices
    drop constraint invoices_customer_id_fkey,
    add constraint invoices_customer_id_fkey foreign key (customer_id) references customers (id);

-- +migrate Down

select 1;