package main

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
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
	"github.com/morning-night-guild/platform-app/internal/driver/server"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
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

	secret := auth.Secret(cfg.JWTSecret)

	authSignUp := interactor.NewAPIAuthSignUp(userRPC, authRPC)

	authSignIn := interactor.NewAPIAuthSignIn(secret, authRPC, authCache, sessionCache)

	authSignOut := interactor.NewAPIAuthSignOut(secret, authCache, sessionCache)

	authVerify := interactor.NewAPIAuthVerify(secret, authCache)

	authRefresh := interactor.NewAPIAuthRefresh(secret, codeCache, authCache, sessionCache)

	authGenerateCode := interactor.NewAPIAuthGenerateCode(secret, codeCache)

	articleList := interactor.NewAPIArticleList(articleRPC)

	articleShare := interactor.NewAPIArticleShare(articleRPC)

	healthUsecase := interactor.NewAPIHealthCheck(healthRPC)

	auth := handler.NewAuth(
		authSignUp,
		authSignIn,
		authSignOut,
		authVerify,
		authRefresh,
		authGenerateCode,
		cookie.New(cfg.CookieDomain),
	)

	article := handler.NewArticle(articleList, articleShare)

	health := handler.NewHealth(healthUsecase)

	hd := http.NewOpenAPI(
		handler.New(cfg.APIKey, secret, auth, article, health),
		cs,
		middleware.New(),
	)

	srv := server.NewServer(cfg.Port, hd)

	srv.Run()
}
