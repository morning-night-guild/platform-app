package main

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/driver/config"
	"github.com/morning-night-guild/platform-app/internal/driver/connect"
	"github.com/morning-night-guild/platform-app/internal/driver/cookie"
	"github.com/morning-night-guild/platform-app/internal/driver/cors"
	"github.com/morning-night-guild/platform-app/internal/driver/env"
	"github.com/morning-night-guild/platform-app/internal/driver/firebase"
	"github.com/morning-night-guild/platform-app/internal/driver/http"
	"github.com/morning-night-guild/platform-app/internal/driver/middleware"
	"github.com/morning-night-guild/platform-app/internal/driver/redis"
	"github.com/morning-night-guild/platform-app/internal/driver/resend"
	"github.com/morning-night-guild/platform-app/internal/driver/server"
)

//nolint:funlen,cyclop
func main() {
	env.Init()

	cfg := config.NewAPI()

	con := connect.New()

	origins, err := cors.ConvertAllowOrigins(cfg.CORSAllowOrigins)
	if err != nil {
		panic(err)
	}

	cs, err := cors.New(origins, cors.ConvertDebugEnable(cfg.CORSDebugEnable))
	if err != nil {
		panic(err)
	}

	noticeRPC := resend.New().MockNotice()
	if cfg.ResendAPIKey != "" {
		noticeRPC, err = resend.New().Notice(cfg.ResendAPIKey, cfg.ResendSender)
		if err != nil {
			panic(err)
		}
	}

	userRPC, err := con.User(cfg.AppCoreURL)
	if err != nil {
		panic(err)
	}

	authRPC, err := firebase.New().Of(cfg.FirebaseSecret, cfg.FirebaseAPIEndpoint, cfg.FirebaseAPIKey)
	if err != nil {
		panic(err)
	}

	articleRPC, err := con.Article(cfg.AppCoreURL)
	if err != nil {
		panic(err)
	}

	healthRPC, err := con.Health(cfg.AppCoreURL)
	if err != nil {
		panic(err)
	}

	rds, err := redis.NewRedis(cfg.RedisURL)
	if err != nil {
		panic(err)
	}

	invitationCache, err := redis.New[model.Invitation]().KVS("invitation", rds)
	if err != nil {
		panic(err)
	}

	userCache, err := redis.New[model.User]().KVS("user", rds)
	if err != nil {
		panic(err)
	}

	authCache, err := redis.New[model.Auth]().KVS("auth", rds)
	if err != nil {
		panic(err)
	}

	sessionCache, err := redis.New[model.Session]().KVS("session", rds)
	if err != nil {
		panic(err)
	}

	codeCache, err := redis.New[model.Code]().KVS("code", rds)
	if err != nil {
		panic(err)
	}

	authUsecase := interactor.NewAPIAuth(
		noticeRPC,
		authRPC,
		userRPC,
		invitationCache,
		userCache,
		authCache,
		codeCache,
		sessionCache,
	)

	articleUsecase := interactor.NewAPIArticle(
		authCache,
		articleRPC,
	)

	healthUsecase := interactor.NewAPIHealth(healthRPC)

	si := handler.New(
		cfg.APIKey,
		auth.Secret(cfg.JWTSecret),
		cookie.New(cfg.CookieDomain),
		authUsecase,
		articleUsecase,
		healthUsecase,
	)

	hd := http.NewOpenAPI(
		si,
		cs,
		middleware.New(),
	)

	srv := server.NewServer(cfg.Port, hd)

	srv.Run()
}
