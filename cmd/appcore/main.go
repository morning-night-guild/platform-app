package main

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/driver/config"
	"github.com/morning-night-guild/platform-app/internal/driver/database"
	"github.com/morning-night-guild/platform-app/internal/driver/env"
	"github.com/morning-night-guild/platform-app/internal/driver/http"
	"github.com/morning-night-guild/platform-app/internal/driver/interceptor"
	"github.com/morning-night-guild/platform-app/internal/driver/newrelic"
	"github.com/morning-night-guild/platform-app/internal/driver/server"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
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

	articleShare := interactor.NewCoreArticleShare(articleRepo)

	articleList := interactor.NewCoreArticleList(articleRepo)

	ctl := controller.New()

	articleCtr := controller.NewArticle(ctl, articleShare, articleList)

	healthCtr := controller.NewHealth()

	var nr *newrelic.NewRelic

	if env.Get().IsProd() {
		nr, err = newrelic.New(cfg.NewRelicAppName, cfg.NewRelicLicense)
		if err != nil {
			panic(err)
		}
	}

	ic := interceptor.New()

	h := http.NewConnect(ic, nr, articleCtr, healthCtr)

	srv := server.NewServer(cfg.Port, h)

	srv.Run()
}
