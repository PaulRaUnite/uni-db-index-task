package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/api"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/config"
)

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
		graceful := make(chan os.Signal, 2)
		signal.Notify(graceful, syscall.SIGTERM)
		signal.Notify(graceful, syscall.SIGINT)

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
		cancel()
		wg.Wait()

	case migrateDownCmd.FullCommand():
		//applied, err := migrate.Exec(cfg.Storage().DB().RawDB(), "postgres", migrations, migrate.Down)
		//if err != nil {
		//	cfg.Log().WithError(err).Fatal("failed to apply migrations")
		//}
		//cfg.Log().WithFields(logan.F{
		//	"direction": "down",
		//	"applied":   applied,
		//}).Info("migrations applied")
	case migrateUpCmd.FullCommand():
		//applied, err := migrate.Exec(cfg.Storage().DB().RawDB(), "postgres", migrations, migrate.Up)
		//if err != nil {
		//	cfg.Log().WithError(err).Fatal("failed to apply migrations")
		//}
		//cfg.Log().WithFields(logan.F{
		//	"direction": "up",
		//	"applied":   applied,
		//}).Info("migrations applied")
	}

}
