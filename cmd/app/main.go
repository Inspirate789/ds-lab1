package main

import (
	"context"
	"fmt"
	"github.com/Inspirate789/ds-lab1/internal/pkg/app"
	"github.com/Inspirate789/ds-lab1/internal/user/repository"
	"github.com/Inspirate789/ds-lab1/internal/user/usecase"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type WebApp interface {
	Start() error
	Shutdown(ctx context.Context) error
}

func startApp(webApp WebApp, config app.Config, logger *slog.Logger) {
	logger.Info(fmt.Sprintf("web app starts at %s with configuration: %+v", config.Web.Host+":"+config.Web.Port, config))

	go func() {
		err := webApp.Start()
		if err != nil {
			panic(err)
		}
	}()
}

func shutdownApp(webApp WebApp, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug("shutdown web app ...")

	const shutdownTimeout = time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	err := webApp.Shutdown(ctx)
	if err != nil {
		panic(err)
	}

	cancel()
	logger.Debug("web app exited")
}

func main() {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "configs/app.yaml", "Config file path")
	pflag.Parse()

	config, err := app.ReadLocalConfig(configPath)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect(config.DB.DriverName, config.DB.ConnectionString)
	if err != nil {
		panic(err)
	}

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(config.Logging.Level)}))
	repo := repository.NewSqlxRepository(db, logger)
	useCase := usecase.New(repo, logger)
	webApp := app.NewFiberApp(config.Web, useCase, logger)

	startApp(webApp, config, logger)
	shutdownApp(webApp, logger)
}
