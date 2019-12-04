package models

import "github.com/shopspring/decimal"

type Good struct {
	ID          int             `jsonapi:"primary,goods"`
	Code        string          `jsonapi:"attr,code"`
	Description string          `jsonapi:"attr,description"`
	Price       decimal.Decimal `jsonapi:"attr,price"`
}

type Invoice struct {
	ID                 int           `jsonapi:"primary,invoice"`
	DestinationCountry string        `jsonapi:"attr,destination_country"`
	CustomerID         int           `jsonapi:"attr,customer_id"`
	Customer           Customer      `jsonapi:"relation,customer"`
	InvoiceParts       []InvoicePart `jsonapi:"relation,invoice_parts"`
}

type InvoicePart struct {
	ID        int      `jsonapi:"primary,invoice_part"`
	GoodID    int      `jsonapi:"attr,good_id"`
	Good      Good     `jsonapi:"relation,good"`
	InvoiceID int      `jsonapi:"attr,invoice_id"`
	Invoice   *Invoice `jsonapi:"attr,invoice"`
}

type Customer struct {
	ID int `jsonapi:"primary,customer"`
}
