package user

import (
	"net/http"
	"strconv"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	rawInvoiceID := chi.URLParam(r, "invoice-id")
	invoiceID, err := strconv.Atoi(rawInvoiceID)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"invoice-id": err})...)
		return
	}
	user, err := handlers.UserQ(r).UserByLogin(username)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	invoice, err := handlers.InvoiceQ(r).InvoiceByID(int64(invoiceID))
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get invoices for user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if invoice == nil || invoice.CustomerID != user.ID {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	country, err := handlers.CountryQ(r).CountryByID(invoice.DestinationCountryID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	parts, goods, err := handlers.InvoiceQ(r).PartsAndGoodsByInvoice(invoice.ID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get parts and goods for invoice")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	err = jsonapi.MarshalPayload(w, models.PopulateInvoice(*invoice, *user, country.ReadableName, parts, goods))
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal invoice response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
