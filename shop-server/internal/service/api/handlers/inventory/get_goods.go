package inventory

import (
	"net/http"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetGoods(w http.ResponseWriter, r *http.Request) {
	goods, err := handlers.GoodQ(r).All(data.GoodSelector{})
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to select goods")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	result := make([]*models.Good, 0, len(goods))
	for _, good := range goods {
		result = append(result, models.PopulateGood(good))
	}
	err = jsonapi.MarshalPayload(w, result)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal goods response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
