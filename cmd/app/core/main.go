package main

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/driver/config"
	"github.com/morning-night-guild/platform-app/internal/driver/env"
	"github.com/morning-night-guild/platform-app/internal/driver/http"
	"github.com/morning-night-guild/platform-app/internal/driver/interceptor"
	"github.com/morning-night-guild/platform-app/internal/driver/newrelic"
	"github.com/morning-night-guild/platform-app/internal/driver/postgres"
	"github.com/morning-night-guild/platform-app/internal/driver/server"
)

func main() {
	env.Init()

	cfg := config.NewCore()

	rdb, err := postgres.New().Of(cfg.DSN)
	if err != nil {
		panic(err)
	}

	articleRepo := gateway.NewArticle(rdb)

	userRepo := gateway.NewUser(rdb)

	articleUsecase := interactor.NewCoreArticle(articleRepo)

	userUsecase := interactor.NewCoreUser(userRepo)

	ctl := controller.New()

	articleCtr := controller.NewArticle(ctl, articleUsecase)

	userCtr := controller.NewUser(ctl, userUsecase)

	healthCtr := controller.NewHealth()

	var nr *newrelic.NewRelic

	if env.Get().IsProd() {
		nr, err = newrelic.New(cfg.NewRelicAppName, cfg.NewRelicLicense)
		if err != nil {
			panic(err)
		}
	}

	ic := interceptor.New()

	h := http.NewConnect(ic, nr, articleCtr, userCtr, healthCtr)

	srv := server.NewServer(cfg.Port, h)

	srv.Run()
}
