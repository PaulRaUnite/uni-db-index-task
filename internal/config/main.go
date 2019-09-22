package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	DB() *pgdb.DB
}
type config struct {
	getter kv.Getter

	database comfig.Once
}

func NewConfig(getter kv.Getter) Config {
	return config{
		getter: getter,
	}
}

func (c config) DB() *pgdb.DB {
	return c.database.Do(func() interface{} {
		var cfg struct {
			Database string
			Password string
		}
		err := figure.Out(&cfg).From(kv.MustGetStringMap(c.getter, "postgres")).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out postgres config"))
		}
		return pgdb.NewDatabaser(c.getter).DB()
	}).(*pgdb.DB)
}
