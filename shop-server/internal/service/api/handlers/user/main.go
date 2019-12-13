package user

import (
	"net/http"
	"strconv"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"

	"github.com/google/jsonapi"
)

func Get(w http.ResponseWriter, r *http.Request) {
	selector := data.CustomerSelector{}
	search := r.URL.Query().Get("filter[name]")
	if search != "" {
		selector.Search = &search
	}
	customers, err := handlers.CustomerQ(r).All(selector)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get users")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	users := make([]*models.User, 0, len(customers))
	for _, customer := range customers {
		users = append(users, models.PopulateUser(customer))
	}
	err = jsonapi.MarshalPayload(w, users)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal users response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

type LogInResponse struct {
	JWT string `jsonapi:"primary,jwt"`
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "user-id")
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"user-id": err})...)
		return
	}
	ss, err := handlers.IssueJWT(userID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to issue JWT")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	err = jsonapi.MarshalPayload(w, &LogInResponse{JWT: ss})
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal jwt response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "user-id")
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"user-id": err})...)
		return
	}
	user, err := handlers.CustomerQ(r).CustomerByID(userID)
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

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "user-id")
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"user-id": err})...)
		return
	}
	rawInvoiceID := chi.URLParam(r, "invoice-id")
	invoiceID, err := strconv.Atoi(rawInvoiceID)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"invoice-id": err})...)
		return
	}
	user, err := handlers.CustomerQ(r).CustomerByID(userID)
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
	if invoice == nil || invoice.CustomerID != userID {
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
