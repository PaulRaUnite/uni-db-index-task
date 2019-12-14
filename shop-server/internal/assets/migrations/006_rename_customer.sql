-- +migrate Up

alter table customers
    rename to users;

-- +migrate Down

alter table users
    rename to customers;