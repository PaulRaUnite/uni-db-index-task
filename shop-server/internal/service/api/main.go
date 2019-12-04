package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/config"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api/handlers/inventory"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, config config.Config) {
	r := chi.NewRouter()

	log := config.Log().WithField("service", "api")
	ape.DefaultMiddlewares(r, log, time.Second)

	r.Use(ape.CtxMiddleware(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "config", config)
	}))

	r.Get("/health", func(writer http.ResponseWriter, r *http.Request) {})
	r.Get("/inventory/goods", inventory.GetGoods)
	r.Get("/inventory/goods/{id}", inventory.GetSingleGood)

	log.Info("started")
	ape.Serve(ctx, r, config, ape.ServeOpts{ShutdownTimeout: 100 * time.Millisecond})
}
