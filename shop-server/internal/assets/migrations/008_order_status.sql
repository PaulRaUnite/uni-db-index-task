-- +migrate Up

alter table invoices
    add column status varchar not null default 'delivered';

-- +migrate Down

alter table invoices
    drop column status;
