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

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "user-id")
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"user-id": err})...)
		return
	}
	user, err := handlers.UserQ(r).UserByID(userID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	invoices, err := handlers.InvoiceQ(r).InvoicesByUser(userID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get invoices for user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := make([]*models.Invoice, 0, len(invoices))
	for _, invoice := range invoices {
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
		result = append(result, models.PopulateInvoice(invoice, *user, country.ReadableName, parts, goods))
	}
	err = jsonapi.MarshalPayload(w, result)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal invoices response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
