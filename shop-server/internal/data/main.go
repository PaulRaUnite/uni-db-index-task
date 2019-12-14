package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
)

//go:generate xo pgsql://postgres:postgres@localhost/shop?sslmode=disable -o ./ -p data --template-path templates
//go:generate xo pgsql://postgres:postgres@localhost/shop?sslmode=disable -o postgres --template-path postgres/templates

type Storage interface {
	Clone() Storage
	DB() *pgdb.DB
	Transaction(tx func() error) error
	GoodQ() GoodQ
	UserQ() UserQ
	InvoiceQ() InvoiceQ
	InvoicePartQ() InvoicePartQ
	CountryQ() CountryQ
}

type GoodSelector struct {
}

type GoodQ interface {
	All(selector GoodSelector) ([]Good, error)
	GoodByID(id int) (*Good, error)
}

type CountryQ interface {
	CountryByID(id int) (*Country, error)
}

type UserSelector struct {
	Name *string
}

type UserQ interface {
	UserByID(id int) (*User, error)
	All(selector UserSelector) ([]User, error)
	UserByLogin(login string) (*User, error)
	Insert(u *User) error
}

type InvoiceQ interface {
	InvoiceByID(id int64) (*Invoice, error)
	InvoicesByUser(user_id int) ([]Invoice, error)
	PartsAndGoodsByInvoice(invoiceID int64) ([]InvoicePart, []Good, error)
}

type InvoicePartQ interface {
	InvoicePartByID(id int64) (*InvoicePart, error)
}
