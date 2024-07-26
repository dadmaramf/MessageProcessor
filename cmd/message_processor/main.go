package main

import (
	"messageprocessor/internal/config"
	"messageprocessor/pkg/postgres"

	"go.uber.org/fx"
)

func main() {

}

func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewConfig,
			postgres.NewPostgresDB,
		),
		// fx.Invoke(

		// )
	)
}
