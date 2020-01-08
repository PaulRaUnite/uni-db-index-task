package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (q InvoiceQ) InvoicesByUser(userID int) ([]data.Invoice, error) {
	query := squirrel.Select("*").From("invoices").Where("customer_id = ?", userID)
	var invoices []data.Invoice
	err := q.db.Select(&invoices, query)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (q InvoiceQ) PartsAndGoodsByInvoice(invoiceID int64) ([]data.InvoicePart, []data.Good, error) {
	query := squirrel.Select("invoice_parts.*, goods.id as good_iid, goods.code as good_code, goods.description as good_descr,goods.price as good_price").
		From("invoice_parts").
		Join("goods on invoice_parts.good_id = goods.id").
		Where("invoice_parts.invoice_id = ?", invoiceID)

	var raw []struct {
		data.InvoicePart
		GoodID    int     `db:"good_iid"`
		GoodCode  string  `db:"good_code"`
		GoodDescr string  `db:"good_descr"`
		GoodPrice float64 `db:"good_price"`
	}
	err := q.db.Select(&raw, query)
	if err != nil {
		return nil, nil, err
	}
	parts := make([]data.InvoicePart, 0, len(raw))
	goods := make([]data.Good, 0, len(raw))
	for _, s := range raw {
		parts = append(parts, s.InvoicePart)
		goods = append(goods, data.Good{
			ID:          s.GoodID,
			Code:        s.GoodCode,
			Description: s.GoodDescr,
			Price:       s.GoodPrice,
		})
	}
	return parts, goods, nil
}

func (q InvoiceQ) CreateInvoice(invoice data.Invoice, parts []data.InvoicePart) error {
	return q.db.Transaction(func() error {
		err := q.Insert(&invoice)
		if err != nil {
			return err
		}
		for _, part := range parts {
			part.InvoiceID = invoice.ID
			err = q.InsertPart(&part)
			if err != nil {
				return err
			}
			good, err := q.GoodByID(part.GoodID)
			if err != nil {
				return err
			}
			if int(good.Amount) < part.Quantity {
				return errors.New("insufficient quantity")
			}
			good.Amount -= int16(part.Quantity)
			err = q.UpdateGood(good)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// InsertPart inserts the InvoicePart to the database.
func (q InvoiceQ) InsertPart(ip *data.InvoicePart) error {
	var err error

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.invoice_parts (` +
		`good_id, unit_price, quantity, invoice_id` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	// run query
	err = q.db.GetRaw(&ip.ID, sqlstr, ip.GoodID, ip.UnitPrice, ip.Quantity, ip.InvoiceID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the Good in the database.
func (q InvoiceQ) UpdateGood(g *data.Good) error {
	var err error

	// sql query
	const sqlstr = `UPDATE public.goods SET (` +
		`code, description, price, amount` +
		`) = ROW( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`
	// run query
	err = q.db.ExecRaw(sqlstr, g.Code, g.Description, g.Price, g.Amount, g.ID)
	return err
}

// GoodByID retrieves a row from 'public.goods' as a Good.
//
// Generated from index 'goods_pkey'.
func (q InvoiceQ) GoodByID(id int) (*data.Good, error) {
	var err error
	// sql query
	const sqlstr = `SELECT ` +
		`id, code, description, price, amount ` +
		`FROM public.goods ` +
		`WHERE id = $1`

	// run query
	g := data.Good{}

	err = q.db.GetRaw(&g, sqlstr, id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &g, nil
}
