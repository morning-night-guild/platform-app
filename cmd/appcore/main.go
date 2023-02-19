package main

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/driver/config"
	"github.com/morning-night-guild/platform-app/internal/driver/database"
	"github.com/morning-night-guild/platform-app/internal/driver/env"
	"github.com/morning-night-guild/platform-app/internal/driver/handler"
	"github.com/morning-night-guild/platform-app/internal/driver/interceptor"
	"github.com/morning-night-guild/platform-app/internal/driver/newrelic"
	"github.com/morning-night-guild/platform-app/internal/driver/server"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor/article"
)

func main() {
	env.Init()

	cfg := config.NewCore()

	db := database.NewClient()

	rdb, err := db.Of(cfg.DSN)
	if err != nil {
		panic(err)
	}

	articleRepo := gateway.NewArticle(rdb)

	articleShareItr := article.NewShareInteractor(articleRepo)

	articleListItr := article.NewListInteractor(articleRepo)

	ctl := controller.New()

	articleCtr := controller.NewArticle(ctl, articleShareItr, articleListItr)

	healthCtr := controller.NewHealth()

	var nr *newrelic.NewRelic

	if env.Get().IsProd() {
		nr, err = newrelic.New(cfg.NewRelicAppName, cfg.NewRelicLicense)
		if err != nil {
			panic(err)
		}
	}

	ic := interceptor.New()

	h := handler.NewConnectHandler(ic, nr, articleCtr, healthCtr)

	srv := server.NewServer(cfg.Port, h)

	srv.Run()
}
