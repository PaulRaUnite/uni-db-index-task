package api

import (
	"context"
	"net/http"
	"time"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers/survey"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers/complaint"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers/user"

	"github.com/go-chi/chi"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/config"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers/inventory"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, config config.Config) {
	r := chi.NewRouter()

	log := config.Log().WithField("service", "api")
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			handler.ServeHTTP(w, r)
		})
	})
	ape.DefaultMiddlewares(r, log, time.Second)

	r.Use(ape.CtxMiddleware(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "config", config)
	}))

	r.Get("/health", func(writer http.ResponseWriter, r *http.Request) {})
	r.Get("/user", user.Get)
	r.Get("/user/{user-id}/login", user.LogIn)
	r.Get("/user/{user-id}/invoice", user.GetInvoices)
	r.Get("/user/{user-id}/invoice/{invoice-id}", user.GetInvoice)
	r.Get("/inventory/good", inventory.GetGoods)
	r.Get("/inventory/good/{id}", inventory.GetSingleGood)
	r.Get("/complaint", complaint.GetAll)
	r.Post("/complaint", complaint.Create)
	r.Get("/complaint/{id}", complaint.Get)
	r.Patch("/complaint/{id}", complaint.Review)
	r.Get("/survey", survey.GetAll)
	r.Post("/survey", survey.Answer)
	r.Get("/survey/{id}", survey.Get)
	r.Get("/survey/template", survey.GetTemplates)
	r.Post("/survey/template", survey.CreateTemplate)
	r.Get("/survey/template/{id}", survey.GetTemplate)
	r.Patch("/survey/template/{id}", survey.UpdateTemplate)

	log.Info("started")
	ape.Serve(ctx, r, config, ape.ServeOpts{ShutdownTimeout: 100 * time.Millisecond})
}
