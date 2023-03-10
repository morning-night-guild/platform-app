package main

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/external"
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/driver/config"
	"github.com/morning-night-guild/platform-app/internal/driver/connect"
	"github.com/morning-night-guild/platform-app/internal/driver/cors"
	"github.com/morning-night-guild/platform-app/internal/driver/env"
	"github.com/morning-night-guild/platform-app/internal/driver/http"
	"github.com/morning-night-guild/platform-app/internal/driver/middleware"
	"github.com/morning-night-guild/platform-app/internal/driver/server"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
)

func main() {
	env.Init()

	cfg := config.NewAPI()

	con, err := connect.New().Of(cfg.AppCoreURL)
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

	articleRPC := external.NewArticle(con)

	healthRPC := external.NewHealth(con)

	articleList := interactor.NewAPIArticleList(articleRPC)

	articleShare := interactor.NewAPIArticleShare(articleRPC)

	healthUsecase := interactor.NewAPIHealthCheck(healthRPC)

	article := handler.NewArticle(articleList, articleShare)

	health := handler.NewHealth(healthUsecase)

	hd := http.NewOpenAPI(
		handler.New(cfg.APIKey, article, health),
		cs,
		middleware.New(),
	)

	srv := server.NewServer(cfg.Port, hd)

	srv.Run()
}
