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


create table goods
(
    description varchar,
    price       float,
    code        varchar primary key
);

create type order_part (
    amount int,
    good_description varchar,
    good_price float,
    good_code varchar
    );

create table orders
(
    customer_id         int,
    invoice_no          varchar,
    destination_country varchar,
    created_at          timestamp,
    ordered             list<frozen<order_part>>,

    primary key ( customer_id, invoice_no )
);

create table orders_raw
(
    destination_country varchar,
    created_at          timestamp,
    stock_code          varchar,
    quantity            int,
    unit_price          float,
    primary key ( stock_code, created_at )
);

create type good_stat (
    price float,
    bought int,
    total_price float
    );

create table orders_stats_by_keyword
(
    keyword      varchar,
    date         date,

    order_amount int,
    total_price  float,
    goods        map<varchar, frozen<good_stat>>,


    primary key (keyword, date)
);

create table orders_stats_by_county
(
    country      varchar,
    date         date,

    order_amount int,
    total_price  float,

    goods        map<varchar, frozen<good_stat>>,
    primary key (country, date)
);

create table orders_stats_by_good
(
    good_code    varchar,
    date         date,

    order_amount int,
    total_price  float,
    goods        map<varchar, frozen<good_stat>>,

    primary key (good_code, date)
);

select *
from orders_stats_by_keyword;

select *
from orders_stats_by_county;

select total_price, order_amount
from orders_stats_by_keyword
where keyword = 'tea'
  and "date" > '2010-01-01'
group by keyword;
