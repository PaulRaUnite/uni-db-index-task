package config

import (
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data/postgres"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Listenerer
	comfig.Logger
	ClonedStorage() data.Storage
}

func New(getter kv.Getter) Config {
	return &config{
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser:  pgdb.NewDatabaser(getter),
	}
}

type config struct {
	comfig.Listenerer
	comfig.Logger
	pgdb.Databaser

	storage comfig.Once
}

func (c *config) ClonedStorage() data.Storage {
	return c.storage.Do(func() interface{} {
		return data.Storage(postgres.New(c.DB()))
	}).(data.Storage)
}
