package inventory

import (
	"net/http"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/urlval"
)

type getCountParams struct {
	Description *string `filter:"description"`
}

func GetGoodsCount(w http.ResponseWriter, r *http.Request) {
	selector := getCountParams{}
	err := urlval.Decode(r.URL.Query(), &selector)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	goodsCount, err := handlers.GoodQ(r).AllCount(selector.Description)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to select goods")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	counter := models.Counter{Value: goodsCount}
	err = jsonapi.MarshalPayload(w, &counter)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal goods response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
