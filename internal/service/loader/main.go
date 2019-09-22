package loader

import (
	"encoding/csv"
	"io"
	"log"

	"github.com/shopspring/decimal"

	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/internal/config"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/PaulRaUnite/uni-db-index-task/internal/data"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

type impl struct {
	db        *pgdb.DB
	source    *gocsv.Unmarshaller
	customers map[int]struct{}
	goods     map[string]int
	invoices  map[int]struct{}
}

func Run(cfg config.Config, source io.Reader) error {
	unm, err := gocsv.NewUnmarshaller(csv.NewReader(source), data.Record{})
	if err != nil {
		return errors.Wrap(err, "failed to create unmarshaller")
	}
	i := impl{
		db:        cfg.DB(),
		source:    unm,
		customers: make(map[int]struct{}, 512),
		goods:     make(map[string]int, 512),
		invoices:  make(map[int]struct{}, 512),
	}
	return i.runOnce()
}

func (i *impl) runOnce() error {
	for {
		record, err := i.source.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal record")
		}
		err = i.upsertInvoicePart(record.(data.Record))
		if err != nil {
			log.Println(record)
			return errors.Wrap(err, "failed to save record in database")
		}
	}
}

func (i *impl) upsertInvoicePart(record data.Record) error {
	err := i.upsertCustomer(record.CustomerID)
	if err != nil {
		return err
	}
	err = i.upsertInvoice(record)
	if err != nil {
		return err
	}
	goodID, err := i.upsertGood(record.StockCode, record.UnitPrice, record.Description)
	if err != nil {
		return err
	}

	return i.insertInvoicePart(goodID, record)
}

func (i *impl) insertInvoicePart(goodID int, record data.Record) error {
	query := squirrel.Insert("invoice_parts").
		SetMap(map[string]interface{}{
			"good_id":    goodID,
			"unit_price": record.UnitPrice,
			"quantity":   record.Quantity,
			"invoice_id": record.InvoiceNo,
		})
	err := i.db.Exec(query)
	return errors.Wrap(err, "failed to insert invoice part")
}

const postgrestimestampFormat = "2006-01-02 15:04:05-07"

func (i *impl) upsertInvoice(record data.Record) error {
	if _, ok := i.invoices[record.InvoiceNo]; ok {
		return nil
	}

	selectQuery := squirrel.Select("id").
		From("invoices").
		Where(squirrel.Eq{"id": record.InvoiceNo})
	var ids []struct{ ID int }
	err := i.db.Select(&ids, selectQuery)
	if err == nil && len(ids) != 0 {
		i.invoices[record.InvoiceNo] = struct{}{}
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to select invoices")
	}

	insertQuery := squirrel.Insert("invoices").
		SetMap(map[string]interface{}{
			"id":                  record.InvoiceNo,
			"customer_id":         record.CustomerID,
			"destination_country": record.Country,
			"invoice_date":        record.InvoiceDate.Format(postgrestimestampFormat),
		})

	err = i.db.Exec(insertQuery)
	if err != nil {
		return errors.Wrap(err, "failed to insert into invoices")
	}
	i.invoices[record.InvoiceNo] = struct{}{}
	return nil
}

func (i *impl) upsertCustomer(customerID int) error {
	if _, ok := i.customers[customerID]; ok {
		return nil
	}

	selectQuery := squirrel.Select("id").
		From("customers").
		Where(squirrel.Eq{"id": customerID})
	var ids []struct{ ID int }
	err := i.db.Select(&ids, selectQuery)
	if err == nil && len(ids) != 0 {
		i.customers[customerID] = struct{}{}
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to select good")
	}

	query := squirrel.Insert("customers").
		SetMap(map[string]interface{}{
			"id": customerID,
		})

	err = i.db.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to insert customer")
	}
	i.customers[customerID] = struct{}{}
	return nil
}

func (i *impl) upsertGood(code string, price decimal.Decimal, description string) (id int, err error) {
	newID, ok := i.goods[code]
	if ok {
		return newID, nil
	}

	selectQuery := squirrel.Select("id").
		From("goods").
		Where(squirrel.Eq{"code": code})
	var ids []struct{ ID int }
	err = i.db.Select(&ids, selectQuery)
	if err == nil && len(ids) != 0 {
		newID = ids[0].ID
		err = i.db.Exec(squirrel.Update("goods").
			Where(squirrel.Eq{"id": newID}).
			SetMap(map[string]interface{}{"price": price}),
		)
		if err != nil {
			return 0, errors.Wrap(err, "failed to update good price")
		}
		i.goods[code] = newID
		return newID, nil
	}
	if err != nil {
		return 0, errors.Wrap(err, "failed to select good")
	}

	insertQuery := squirrel.Insert("goods").
		SetMap(map[string]interface{}{
			"code":        code,
			"price":       price,
			"description": description,
		}).Suffix("RETURNING id")

	err = i.db.Select(&ids, insertQuery)
	if err != nil {
		log.Println(description)
		return 0, errors.Wrap(err, "failed to insert into goods")
	}
	newID = ids[0].ID
	i.goods[code] = newID
	return newID, nil
}
