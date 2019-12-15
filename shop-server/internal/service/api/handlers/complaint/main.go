package complaint

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/jsonapi"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	cursor, err := handlers.ComplaintsQ(r).Find(r.Context(), bson.D{})
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
	if err = cursor.Err(); err != nil {
		ape.Log(r).WithError(err).Error("failed to get complaints")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	err = jsonapi.MarshalPayload(w, result)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal complaints payload")
		ape.RenderErr(w, problems.InternalError())
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer io.Copy(ioutil.Discard, r.Body)

	userID, err := handlers.UserIDFromClaims(r)
	if err != nil {
		ape.Log(r).WithError(err).Warn("failed to get user id")
		ape.RenderErr(w, problems.NotAllowed(err))
		return
	}
	user, err := handlers.UserQ(r).UserByID(userID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get customer")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	complaint := models.Complaint{}
	err = jsonapi.UnmarshalPayload(r.Body, &complaint)
	if err != nil {
		ape.Log(r).WithError(err).Warn("bad request")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"body": err})...)
		return
	}
	complaint.ID_ = primitive.NewObjectID()
	complaint.Answer = ""
	complaint.Reviewer = nil
	complaint.User = models.PopulateUser(*user)
	complaint.CreatedAt = time.Now()

	_, err = handlers.ComplaintsQ(r).InsertOne(r.Context(), &complaint)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to insert complaint")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	complaint.ID = complaint.ID_.Hex()

	err = jsonapi.MarshalPayload(w, &complaint)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to render complaint")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	objID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"id": err})...)
		return
	}
	complaint := models.Complaint{}
	err = handlers.ComplaintsQ(r).FindOne(r.Context(), bson.M{"_id": objID}).Decode(&complaint)
	if err != nil {
		if errors.Cause(err) == mongo.ErrNoDocuments {
			ape.RenderErr(w, problems.NotFound())
			return
		}
		ape.Log(r).WithError(err).Error("failed to get complaints cursor")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	complaint.ID = complaint.ID_.Hex()
	err = jsonapi.MarshalPayload(w, &complaint)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal complaints payload")
		ape.RenderErr(w, problems.InternalError())
	}
}

func Review(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer io.Copy(ioutil.Discard, r.Body)

	complaint := models.Complaint{}
	err := jsonapi.UnmarshalPayload(r.Body, &complaint)
	if err != nil {
		ape.Log(r).WithError(err).Warn("bad request")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"body": err})...)
		return
	}
	userID, err := handlers.UserIDFromClaims(r)
	if err != nil {
		ape.Log(r).WithError(err).Warn("failed to get user id")
		ape.RenderErr(w, problems.NotAllowed(err))
		return
	}
	user, err := handlers.UserQ(r).UserByID(userID)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get customer")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	objID, err := primitive.ObjectIDFromHex(complaint.ID)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"id": err})...)
		return
	}
	result, err := handlers.ComplaintsQ(r).UpdateOne(r.Context(), bson.M{"_id": objID},
		bson.M{"$set": bson.M{"reviewer": user, "answer": complaint.Answer, "reviewed_at": time.Now()}})
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to update complaint")
		ape.RenderErr(w, problems.InternalError())
	}
	if result.ModifiedCount == 0 {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	Get(w, r)
}
