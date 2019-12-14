package user

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Get(w http.ResponseWriter, r *http.Request) {
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
	claims, err := handlers.GetClaims(r)
	if err != nil || claims == nil || claims.UserID != user.ID {
		ape.RenderErr(w, problems.NotAllowed())
		return
	}

	err = jsonapi.MarshalPayload(w, models.PopulateUser(*user))
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal invoices response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
