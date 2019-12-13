package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	migrate "github.com/rubenv/sql-migrate"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/assets"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/config"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/service/api"
)

var migrations = &migrate.PackrMigrationSource{Box: assets.SchemaMigrations}

func main() {
	defer func() {
		if rvr := recover(); rvr != nil {
			logan.New().WithRecover(rvr).Fatal("app panicked")
		}
	}()

	app := kingpin.New("leasing-svc", "")
	runCmd := app.Command("run", "")
	migrateCmd := app.Command("migrate", "")
	migrateUpCmd := migrateCmd.Command("up", "")
	migrateDownCmd := migrateCmd.Command("down", "")

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		logan.New().WithError(err).Fatal("failed to parse arguments")
	}

	cfg := config.New(kv.MustFromEnv())
	switch cmd {
	case runCmd.FullCommand():
		ctx, cancel := context.WithCancel(context.Background())
		graceful := make(chan os.Signal, 3)
		signal.Notify(graceful, syscall.SIGTERM)
		signal.Notify(graceful, syscall.SIGINT)
		signal.Notify(graceful, syscall.SIGKILL)

		routines := []func(context.Context, config.Config){
			api.Run,
		}

		wg := sync.WaitGroup{}
		wg.Add(len(routines))
		for _, f := range routines {
			go func(f func(context.Context, config.Config)) {
				f(ctx, cfg)
				wg.Done()
			}(f)
		}

		<-graceful
		cfg.Log().Info("graceful shutdown")
		cancel()
		wg.Wait()
	case migrateDownCmd.FullCommand():
		applied, err := migrate.Exec(cfg.ClonedStorage().DB().RawDB(), "postgres", migrations, migrate.Down)
		if err != nil {
			cfg.Log().WithError(err).Fatal("failed to apply migrations")
		}
		cfg.Log().WithFields(logan.F{
			"direction": "down",
			"applied":   applied,
		}).Info("migrations applied")
	case migrateUpCmd.FullCommand():
		applied, err := migrate.Exec(cfg.ClonedStorage().DB().RawDB(), "postgres", migrations, migrate.Up)
		if err != nil {
			cfg.Log().WithError(err).Fatal("failed to apply migrations")
		}
		cfg.Log().WithFields(logan.F{
			"direction": "up",
			"applied":   applied,
		}).Info("migrations applied")
	}

}
