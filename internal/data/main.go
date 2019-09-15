package data

import (
	"time"

	"github.com/shopspring/decimal"
)

type Time struct {
	time.Time
}

const format = "1/2/2006 15:04"

func (t *Time) UnmarshalText(data []byte) error {
	tt, err := time.Parse(format, string(data))
	if err != nil {
		return err
	}
	*t = Time{Time: tt}
	return nil
}

type Record struct {
	InvoiceNo   int
	StockCode   string
	Description string
	Quantity    int
	InvoiceDate Time
	UnitPrice   decimal.Decimal
	CustomerID  int
	Country     string
}
