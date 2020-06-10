/*
Categorization/filtering parameters:
- defined good
- good category
- country of order
- year or period of time

Target paramaters of interest:
- average cart size
- average cart price
- busyest days/months in the defined period
- sale dynamics (current/last year or similar)
*/

CREATE KEYSPACE shop
            WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 3};

use shop;

create table users_by_id
(
    id        bigint primary key,
    username  varchar,
    pass_hash varchar
);

create table users_by_credentials
(
    id        bigint,
    username  varchar,
    pass_hash varchar,
    primary key ((username, pass_hash))
);

create table goods_by_id
(
    description varchar,
    price       float,
    code        varchar,
    id          bigint primary key,
);

create table goods_by_keyword_asc
(
    keyword     varchar,
    description varchar,
    price       float,
    code        varchar,
    id          bigint,
    primary key (keyword, price, id )
) with clustering order by (price asc);

create table goods_by_keyword_desc
(
    keyword varchar,
    price   float,
    code    varchar,
    id      bigint,
    primary key (keyword, price, id )
) with clustering order by (price desc);

create type order_part (
    amount int,
    good_description varchar,
    good_price float,
    good_code varchar,
    good_id bigint,
    );

create table orders_by_customer
(
    id                  bigint,
    customer_id         int,
    destination_country varchar,
    created_at          timestamp,
    ordered             list<frozen<order_part>>,
    primary key (customer_id, created_at)
) with clustering order by (created_at asc);

create table orders_by_year_month
(
    id                  bigint,
    customer_id         int,
    destination_country varchar,
    created_at          timestamp,
    ordered             list<frozen<order_part>>,

    year                int,
    month               int,
    day                 int,
    primary key ((year, month), day, id)
) with clustering order by (day asc, id asc);

create table orders_by_country_year_month
(
    id                  bigint,
    customer_id         int,
    destination_country varchar,
    created_at          timestamp,
    ordered             list<frozen<order_part>>,

    year                int,
    month               int,
    day                 int,
    primary key ((destination_country, year, month), day, id)
) with clustering order by (day asc, id asc);

select *
from goods_by_id
where id = 1;

select *
from goods_by_keyword_asc
where keyword in ('tea', 'biscuits');

select *
from users_by_credentials
where username = 'merchant_from_france'
  and pass_hash = 'asa2fa12tu78osdhas512hd7i';

select *
from orders_by_customer
where customer_id = 1;

select *
from orders_by_year_month
where year = 2020
  and month in (1,2,3,4);

select *
from orders_by_country_year_month
where destination_country = 'france'
  and year = 2020
  and month in (1,2,3,4);
