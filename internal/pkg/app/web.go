package app

import (
	"context"
	"github.com/Inspirate789/ds-lab1/internal/person/delivery"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pkg/errors"
	"github.com/samber/slog-fiber"
	"log/slog"
	"net/http"
)

type HealthChecker interface {
	HealthCheck() error
}

type WebConfig struct {
	Host       string `koanf:"host"`
	Port       string `koanf:"port"`
	PathPrefix string `koanf:"path_prefix"`
}

type FiberApp struct {
	config WebConfig
	fiber  *fiber.App
	logger *slog.Logger
}

func newFiberError(msg string) map[string]any {
	return fiber.Map{"message": msg}
}

func checkLiveness(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("live")
}

func checkReadiness(useCase HealthChecker) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		err := useCase.HealthCheck()
		if err != nil {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(newFiberError(err.Error()))
		}

		return ctx.Status(fiber.StatusOK).SendString("ready")
	}
}

func NewFiberApp(config WebConfig, useCase delivery.UseCase, logger *slog.Logger) *FiberApp {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).JSON(newFiberError(err.Error()))
		},
	})

	app.Use(recover.New())
	app.Use(slogfiber.New(logger))
	app.Use(pprof.New())

	app.Get("/health/live", checkLiveness)
	app.Get("/health/ready", checkReadiness(useCase))

	api := app.Group(config.PathPrefix)
	delivery.AddHandlers(api.Group("/persons"), useCase, logger)

	return &FiberApp{
		config: config,
		fiber:  app,
		logger: logger,
	}
}

func (f *FiberApp) Start() error {
	return errors.Wrap(f.fiber.Listen(f.config.Host+":"+f.config.Port), "start web app")
}

func (f *FiberApp) Shutdown(ctx context.Context) error {
	return errors.Wrap(f.fiber.ShutdownWithContext(ctx), "stop web app")
}

func (f *FiberApp) Test(req *http.Request, msTimeout ...int) (*http.Response, error) {
	return f.fiber.Test(req, msTimeout...)
}
