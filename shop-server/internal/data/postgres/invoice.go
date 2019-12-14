package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
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
