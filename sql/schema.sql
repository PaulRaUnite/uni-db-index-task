create table Customers (
    id SERIAL PRIMARY KEY
);

create table Goods (
    id SERIAL PRIMARY KEY,
    description varchar(1024) NOT NULL,
    price REAl
);

create table Invoices (
    id BIGSERIAL PRIMARY KEY,
    stock_code INT REFERENCES Goods(id) NOT NULL,
    invoice_date TIMESTAMP NOT NULL,
    unit_price DECIMAL NOT NULL,
    customer_id INT REFERENCES Customers(id),
    destination_country VARCHAR(128)
)
