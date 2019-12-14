package user

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type LogInResponse struct {
	JWT string `jsonapi:"primary,jwt"`
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"auth": errors.New("no credentials")})...)
		return
	}
	user, err := handlers.UserQ(r).UserByLogin(username)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == nil || user.Password != password {
		ape.RenderErr(w, problems.NotAllowed())
		return
	}
	ss, err := handlers.IssueJWT(*user)
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
