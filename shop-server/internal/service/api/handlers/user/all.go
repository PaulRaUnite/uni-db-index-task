package user

import (
	"net/http"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func All(w http.ResponseWriter, r *http.Request) {
	selector := data.UserSelector{}
	search := r.URL.Query().Get("filter[name]")
	if search != "" {
		selector.Name = &search
	}
	customers, err := handlers.UserQ(r).All(selector)
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
