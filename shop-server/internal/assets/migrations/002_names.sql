-- +migrate Up

alter table customers add column name varchar(256) not null default '';

-- +migrate Down

alter table customers drop column name;
