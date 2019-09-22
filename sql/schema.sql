create table customers
(
    id SERIAL PRIMARY KEY
);

create table goods
(
    id          SERIAL PRIMARY KEY,
    code        VARCHAR(128)  NOT NULL UNIQUE,
    description TEXT NOT NULL ,
    price       DECIMAL       NOT NULL
);

create table invoices
(
    id                  BIGSERIAL PRIMARY KEY,
    customer_id         INT REFERENCES customers (id) NOT NULL,
    destination_country VARCHAR(128)                  NOT NULL,
    invoice_date        TIMESTAMP                     NOT NULL
);

create table invoice_parts
(
    id         BIGSERIAL PRIMARY KEY,
    good_id    INT REFERENCES goods (id)       NOT NULL,
    unit_price DECIMAL                         NOT NULL,
    quantity   INT                             NOT NULL,
    invoice_id BIGINT REFERENCES invoices (id) NOT NULL
);

