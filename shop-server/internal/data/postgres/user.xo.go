// Package postgres contains the types for schema 'public'.
package postgres

import (
	"database/sql"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

// Code generated by xo. DO NOT EDIT.

// UserQ represents helper struct to access row of 'public.users'.
type UserQ struct {
	db *pgdb.DB
}

// NewUserQ  - creates new instance
func NewUserQ(db *pgdb.DB) UserQ {
	return UserQ{
		db,
	}
}

// UserQ  - creates new instance of UserQ
func (s Storage) UserQ() data.UserQ {
	return NewUserQ(s.DB())
}

// Insert inserts the User to the database.
func (q UserQ) Insert(u *data.User) error {
	var err error

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.users (` +
		`name, password, account_type, login` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	// run query
	err = q.db.GetRaw(&u.ID, sqlstr, u.Name, u.Password, u.AccountType, u.Login)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the User in the database.
func (q UserQ) Update(u *data.User) error {
	var err error

	// sql query
	const sqlstr = `UPDATE public.users SET (` +
		`name, password, account_type, login` +
		`) = ROW( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`
	// run query
	err = q.db.ExecRaw(sqlstr, u.Name, u.Password, u.AccountType, u.Login, u.ID)
	return err
}

// Upsert performs an upsert for User.
func (q UserQ) Upsert(u *data.User) error {
	var err error
	// sql query
	const sqlstr = `INSERT INTO public.users (` +
		`id, name, password, account_type, login` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, password, account_type, login` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.password, EXCLUDED.account_type, EXCLUDED.login` +
		`)`
	// run query
	err = q.db.ExecRaw(sqlstr, u.ID, u.Name, u.Password, u.AccountType, u.Login)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes the User from the database.
func (q UserQ) Delete(id int) error {
	var err error
	// sql query with composite primary key
	const sqlstr = `DELETE FROM public.users  WHERE id = $1`

	// run query
	err = q.db.ExecRaw(sqlstr, id)
	if err != nil {
		return err
	}

	return nil
}

// UserByID retrieves a row from 'public.users' as a User.
//
// Generated from index 'customers_pkey'.
func (q UserQ) UserByID(id int) (*data.User, error) {
	var err error
	// sql query
	const sqlstr = `SELECT ` +
		`id, name, password, account_type, login ` +
		`FROM public.users ` +
		`WHERE id = $1`

	// run query
	u := data.User{}

	err = q.db.GetRaw(&u, sqlstr, id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

// UpsertUserByLogin - tries to insert u into db on conflict updates non primary key fields
// and sets generated field into u
func (q UserQ) UpsertUserByLogin(u *data.User) error {
	var err error
	// sql query
	const sqlstr = `INSERT INTO public.users (` +
		`name, password, account_type, login` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (login) DO UPDATE SET ` +
		`(name, password, account_type, login` +
		`) = (` +
		`EXCLUDED.name, EXCLUDED.password, EXCLUDED.account_type, EXCLUDED.login)` +
		` RETURNING id`
	// run query
	err = q.db.GetRaw(&u.ID, sqlstr, u.Name, u.Password, u.AccountType, u.Login)
	if err != nil {
		return err
	}
	return nil
}

// UserByLogin retrieves a row from 'public.users' as a User.
//
// Generated from index 'unique_login'.
func (q UserQ) UserByLogin(login string) (*data.User, error) {
	var err error
	// sql query
	const sqlstr = `SELECT ` +
		`id, name, password, account_type, login ` +
		`FROM public.users ` +
		`WHERE login = $1`

	// run query
	u := data.User{}

	err = q.db.GetRaw(&u, sqlstr, login)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
