package user

import (
	"context"
	"net/http"
	"sort"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"go.mongodb.org/mongo-driver/bson"
)

func GetComplaints(w http.ResponseWriter, r *http.Request) {
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

	cursor, err := handlers.ComplaintsQ(r).Find(r.Context(), bson.M{"user.id": user.ID})
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get complaints cursor")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	defer cursor.Close(context.Background())
	var result []*models.Complaint
	for cursor.Next(context.Background()) {
		complaint := &models.Complaint{}
		err = cursor.Decode(complaint)
		if err != nil {
			ape.Log(r).WithError(err).Error("failed decode complaints")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		complaint.ID = complaint.ID_.Hex()
		result = append(result, complaint)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.After(result[j].Date)
	})
	if err = cursor.Err(); err != nil {
		ape.Log(r).WithError(err).Error("failed to get complaints")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	err = jsonapi.MarshalPayload(w, result)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal complaints payload")
		ape.RenderErr(w, problems.InternalError())
	}
}
