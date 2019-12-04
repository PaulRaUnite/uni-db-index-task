package inventory

import (
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/google/jsonapi"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetSingleGood(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"id": err})...)
		return
	}
	good, err := handlers.GoodQ(r).GoodByID(id)
	if err != nil {
		handlers.Log(r).WithError(err).Error("failed to select good by id")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if good == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	err = jsonapi.MarshalPayload(w, &models.Good{
		ID:          good.ID,
		Code:        good.Code,
		Description: good.Description,
		Price:       decimal.NewFromFloat(good.Price),
	})
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal payload")
		ape.RenderErr(w, problems.InternalError())
	}
}
