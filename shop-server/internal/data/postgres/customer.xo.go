// Package postgres contains the types for schema 'public'.
package postgres

import (
	"database/sql"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

// Code generated by xo. DO NOT EDIT.

// CustomerQ represents helper struct to access row of 'public.customers'.
type CustomerQ struct {
	db *pgdb.DB
}

// NewCustomerQ  - creates new instance
func NewCustomerQ(db *pgdb.DB) CustomerQ {
	return CustomerQ{
		db,
	}
}

// CustomerQ  - creates new instance of CustomerQ
func (s Storage) CustomerQ() data.CustomerQ {
	return NewCustomerQ(s.DB())
}

// Insert inserts the Customer to the database.
func (q CustomerQ) Insert(c *data.Customer) error {
	var err error

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.customers (` +
		`` +
		`) VALUES (` +
		`` +
		`) RETURNING id`

	// run query
	err = q.db.GetRaw(&c.ID, sqlstr)
	if err != nil {
		return err
	}

	return nil
}

// Update statements omitted due to lack of fields other than primary key

// Delete deletes the Customer from the database.
func (q CustomerQ) Delete(id int) error {
	var err error
	// sql query with composite primary key
	const sqlstr = `DELETE FROM public.customers  WHERE id = $1`

	// run query
	err = q.db.ExecRaw(sqlstr, id)
	if err != nil {
		return err
	}

	return nil
}

// CustomerByID retrieves a row from 'public.customers' as a Customer.
//
// Generated from index 'customers_pkey'.
func (q CustomerQ) CustomerByID(id int) (*data.Customer, error) {
	var err error
	// sql query
	const sqlstr = `SELECT ` +
		`id ` +
		`FROM public.customers ` +
		`WHERE id = $1`

	// run query
	c := data.Customer{}

	err = q.db.GetRaw(&c, sqlstr, id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
