package user

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer io.Copy(ioutil.Discard, r.Body)

	userID, err := handlers.UserIDFromClaims(r)
	if err != nil {
		ape.Log(r).WithError(err).Warn("failed to get user id")
		ape.RenderErr(w, problems.NotAllowed(err))
		return
	}
	user, err := handlers.UserQ(r).UserByID(userID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get customer")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	invoice := models.Invoice{}
	err = jsonapi.UnmarshalPayload(r.Body, &invoice)
	if err != nil {
		ape.Log(r).WithError(err).Warn("bad request")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"body": err})...)
		return
	}
	if len(invoice.InvoiceParts) == 0 {
		ape.RenderErr(w, problems.BadRequest(errors.New("no invoice parts"))...)
		return
	}
	country, err := handlers.CountryQ(r).CountryByReadableName(invoice.DestinationCountry)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get country")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if country == nil {
		ape.RenderErr(w, problems.BadRequest(errors.New("no such country"))...)
		return
	}

	invoiceData := data.Invoice{
		CustomerID:           user.ID,
		DestinationCountryID: country.ID,
		InvoiceDate:          time.Now(),
		Status:               "created",
	}
	parts := make([]data.InvoicePart, 0, len(invoice.InvoiceParts))
	for _, part := range invoice.InvoiceParts {
		good, err := handlers.GoodQ(r).GoodByID(part.GoodID)
		if err != nil || good == nil {
			ape.Log(r).Error(err)
			ape.RenderErr(w, problems.InternalError())
			return
		}
		parts = append(parts, data.InvoicePart{
			UnitPrice: good.Price,
			GoodID:    part.GoodID,
			Quantity:  part.Quantity,
		})
	}
	err = handlers.InvoiceQ(r).CreateInvoice(invoiceData, parts)
	if err != nil {
		ape.Log(r).Error(err)
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
