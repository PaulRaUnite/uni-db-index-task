package country

import (
	"net/http"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	countries, err := handlers.CountryQ(r).All()
	if err != nil {
		ape.Log(r).WithError(err).Error("db failed to select countries")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	result := make([]*models.Country, 0, len(countries))
	for _, country := range countries {
		result = append(result, models.PopulateCountry(country))
	}
	err = jsonapi.MarshalPayload(w, result)
	if err != nil {
		ape.Log(r).WithError(err).Error("handler failed to marshal country response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
