package user

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/lib/pq"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/models"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type NewUser struct {
	Login    string `jsonapi:"primary,users"`
	Password string `jsonapi:"attr,password,required"`
	Name     string `jsonapi:"attr,name,required"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer io.Copy(ioutil.Discard, r.Body)

	req := NewUser{}
	err := jsonapi.UnmarshalPayload(r.Body, &req)

	if err != nil {
		ape.Log(r).WithError(err).Error("failed to parse request payload")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = handlers.UserQ(r).Insert(&data.User{
		Name:        req.Name,
		Password:    req.Password,
		AccountType: 0,
		Login:       req.Login,
	})
	if err != nil {
		if e, ok := errors.Cause(err).(*pq.Error); ok && e.Code == "23505" {
			ape.RenderErr(w, problems.Conflict())
			return
		}
		ape.Log(r).WithError(err).Error("failed to insert new user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	user, err := handlers.UserQ(r).UserByLogin(req.Login)
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to get new user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	err = jsonapi.MarshalPayload(w, models.PopulateUser(*user))
	if err != nil {
		ape.Log(r).WithError(err).Error("failed to marshal user response")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
