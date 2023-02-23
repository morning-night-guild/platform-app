package main

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/api"
	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/driver/config"
	"github.com/morning-night-guild/platform-app/internal/driver/connect"
	"github.com/morning-night-guild/platform-app/internal/driver/cors"
	"github.com/morning-night-guild/platform-app/internal/driver/env"
	"github.com/morning-night-guild/platform-app/internal/driver/handler"
	"github.com/morning-night-guild/platform-app/internal/driver/middleware"
	"github.com/morning-night-guild/platform-app/internal/driver/server"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
)

func main() {
	env.Init()

	cfg := config.NewAPI()

	c, err := connect.New().Of(cfg.AppCoreURL)
	if err != nil {
		panic(err)
	}

	origins, err := cors.ConvertAllowOrigins(cfg.CORSAllowOrigins)
	if err != nil {
		panic(err)
	}

	cs, err := cors.New(origins, cors.ConvertDebugEnable(cfg.CORSDebugEnable))
	if err != nil {
		panic(err)
	}

	articleGateway := gateway.NewCoreArticle(c)

	healthGateway := gateway.NewCoreHealth(c)

	articleListUsecase := interactor.NewAPIArticleList(articleGateway)

	articleShareUsecase := interactor.NewAPIArticleShare(articleGateway)

	healthUsecase := interactor.NewAPIHealthCheck(healthGateway)

	article := api.NewArticle(articleListUsecase, articleShareUsecase)

	health := api.NewHealth(healthUsecase)

	hd := handler.NewOpenAPIHandler(
		api.New(cfg.APIKey, article, health),
		cs,
		middleware.New(),
	)

	srv := server.NewServer(cfg.Port, hd)

	srv.Run()
}
