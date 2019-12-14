package user

import (
	"net/http"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
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
	invoices, err := handlers.InvoiceQ(r).InvoicesByUser(user.ID)
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
