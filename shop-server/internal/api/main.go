package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/config"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, config config.Config) {
	r := chi.NewRouter()

	log := config.Log().WithField("service", "api")
	ape.DefaultMiddlewares(r, log, time.Second)

	r.Get("/health", func(writer http.ResponseWriter, r *http.Request) {})

	log.Info("started")
	ape.Serve(ctx, r, config, ape.ServeOpts{ShutdownTimeout: time.Second})
}
