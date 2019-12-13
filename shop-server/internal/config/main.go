package config

import (
	"context"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data/postgres"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config interface {
	comfig.Listenerer
	comfig.Logger
	ClonedStorage() data.Storage
	Surveys() *mongo.Collection
	Complaints() *mongo.Collection
}

func New(getter kv.Getter) Config {

	return &config{
		getter:     getter,
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser:  pgdb.NewDatabaser(getter),
	}
}

type config struct {
	getter kv.Getter
	comfig.Listenerer
	comfig.Logger
	pgdb.Databaser

	storage comfig.Once
	mongo   comfig.Once
}

func (c *config) ClonedStorage() data.Storage {
	return c.storage.Do(func() interface{} {
		return data.Storage(postgres.New(c.DB()))
	}).(data.Storage)
}

func (c *config) MongoStorage() *mongo.Database {
	return c.mongo.Do(func() interface{} {
		var cfg struct {
			URL string `fig:"url,required"`
		}
		err := figure.Out(&cfg).From(kv.MustGetStringMap(c.getter, "mongo")).Please()
		if err != nil {
			panic(err)
		}
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URL).SetAuth(options.Credential{
			Username:    "root",
			Password:    "example",
			PasswordSet: true,
		}))
		if err != nil {
			panic(err)
		}
		return client.Database("junk")
	}).(*mongo.Database)
}

func (c *config) Surveys() *mongo.Collection {
	return c.MongoStorage().Collection("surveys")
}

func (c *config) Complaints() *mongo.Collection {
	return c.MongoStorage().Collection("complaints")
}
