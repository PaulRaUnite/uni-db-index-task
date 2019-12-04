package handlers

import (
	"net/http"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/config"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

func GoodQ(r *http.Request) data.GoodQ {
	return Config(r).ClonedStorage().GoodQ()
}

func Config(r *http.Request) config.Config {
	return r.Context().Value("config").(config.Config)
}

func Log(r *http.Request) *logan.Entry {
	return ape.Log(r)
}
