-- +migrate Up

create extension pg_trgm;

-- +migrate Down

drop extension pg_trgm;
