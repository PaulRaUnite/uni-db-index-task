-- +migrate Up

alter table goods
    add column amount smallint not null default 0;

update goods
set amount = floor(random() * 10 + 1)::int;

-- +migrate Down

alter table goods
    drop column amount;
