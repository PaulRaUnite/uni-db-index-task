package models

import (
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Good struct {
	ID          int             `jsonapi:"primary,goods"`
	Code        string          `jsonapi:"attr,code"`
	Description string          `jsonapi:"attr,description"`
	Price       decimal.Decimal `jsonapi:"attr,price"`
}

type Invoice struct {
	ID                 int64          `jsonapi:"primary,invoices"`
	DestinationCountry string         `jsonapi:"attr,destination_country"`
	CustomerID         int            `jsonapi:"attr,customer_id"`
	Customer           *User          `jsonapi:"relation,customer"`
	InvoiceParts       []*InvoicePart `jsonapi:"relation,invoice_parts"`
}

type InvoicePart struct {
	ID        int64    `jsonapi:"primary,invoice_parts"`
	GoodID    int      `jsonapi:"attr,good_id"`
	Good      *Good    `jsonapi:"relation,good"`
	InvoiceID int64    `jsonapi:"attr,invoice_id"`
	Invoice   *Invoice `jsonapi:"relation,invoice"`
}

type User struct {
	ID   int    `jsonapi:"primary,customers"`
	Name string `jsonapi:"attr,name"`
}

func PopulateGood(good data.Good) *Good {
	return &Good{
		ID:          good.ID,
		Code:        good.Code,
		Description: good.Description,
		Price:       decimal.NewFromFloat(good.Price),
	}
}

func PopulateUser(customer data.Customer) *User {
	return &User{
		ID:   customer.ID,
		Name: customer.Name,
	}
}

func PopulateInvoice(invoice data.Invoice, user data.Customer, country string, parts []data.InvoicePart, goods []data.Good) *Invoice {
	modelParts := make([]*InvoicePart, 0, len(parts))
	for i, part := range parts {
		modelParts = append(modelParts, PopulateInvoicePart(part, goods[i]))
	}
	return &Invoice{
		ID:                 invoice.ID,
		DestinationCountry: country,
		CustomerID:         invoice.CustomerID,
		Customer:           PopulateUser(user),
		InvoiceParts:       modelParts,
	}
}

func PopulateInvoicePart(part data.InvoicePart, good data.Good) *InvoicePart {
	return &InvoicePart{
		ID:        part.ID,
		GoodID:    part.GoodID,
		Good:      PopulateGood(good),
		InvoiceID: part.InvoiceID,
		Invoice:   nil,
	}
}

type Complaint struct {
	ID       string             `bson:"-" jsonapi:"primary,complaints"`
	ID_      primitive.ObjectID `bson:"_id"`
	User     *User              `jsonapi:"relation,user"`
	Body     string             `jsonapi:"attr,body"`
	Answer   string             `jsonapi:"attr,answer"`
	Reviewer *User              `jsonapi:"relation,reviewer"`
}
