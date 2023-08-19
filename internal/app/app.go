package app

import (
	"context"
	"log/slog"

	"github.com/de-bait/internal/db"
	"github.com/de-bait/internal/db/repository"
	"github.com/de-bait/internal/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

type ResourceRegisterer interface {
	Register()
}

func AsResource(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ResourceRegisterer)),
		fx.ResultTags(`group:"resources"`),
	)
}

func NewFiberApp() *fiber.App {
	app := fiber.New()

	return app
}

func NewServerHandler(resources []ResourceRegisterer) func() {
	return func() {
		for _, resource := range resources {
			resource.Register()
		}

	}
}

func NewContext() context.Context {
	return context.Background()
}

func NewValidator() *validator.Validate {
	validate := validator.New()
	return validate
}

func InitApp(lc fx.Lifecycle, app *fiber.App, registerHandlers func()) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			registerHandlers()
			slog.Info("App started on port: :2020")
			go app.Listen(":2020")
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})
}

func Run() {
	fx.New(
		fx.Provide(
			NewFiberApp,
			NewContext,
			NewValidator,
			repository.NewUserRepository,
			db.NewEntClient,
			fx.Annotate(
				NewServerHandler,
				fx.ParamTags(`group:"resources"`),
			),

			AsResource(user.NewUserHandler),
		),
		fx.Invoke(InitApp),
	).Run()

	// Run the auto migration tool.

	// if err := db.Client.Schema.Create(context.Background()); err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }

}
