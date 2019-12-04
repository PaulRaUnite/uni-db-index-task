package data

import "gitlab.com/distributed_lab/kit/pgdb"

//go:generate xo pgsql://postgres:postgres@localhost/shop?sslmode=disable -o ./ -p data --template-path templates
//go:generate xo pgsql://postgres:postgres@localhost/shop?sslmode=disable -o postgres --template-path postgres/templates

type Storage interface {
	Clone() Storage
	DB() *pgdb.DB
	Transaction(tx func() error) error
	GoodQ() GoodQ
}

type GoodQ interface {
	GoodByID(id int) (*Good, error)
}

type CountryQ interface {
}

type CustomerQ interface {
}

type InvoiceQ interface {
}

type InvoicePartQ interface {
}
