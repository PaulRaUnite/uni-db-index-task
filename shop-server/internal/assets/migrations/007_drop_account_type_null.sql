-- +migrate Up

alter table users
    alter column account_type set not null;

-- +migrate Down

alter table users
    alter column account_type drop not null;
