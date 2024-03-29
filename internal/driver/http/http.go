package http

import (
	"net/http"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"github.com/go-chi/chi/v5"
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/driver/middleware"
	"github.com/morning-night-guild/platform-app/internal/driver/newrelic"
	"github.com/morning-night-guild/platform-app/internal/driver/router"
	"github.com/morning-night-guild/platform-app/pkg/connect/article/v1/articlev1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
	"github.com/morning-night-guild/platform-app/pkg/connect/user/v1/userv1connect"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func NewConnect(
	interceptor connect.UnaryInterceptorFunc,
	nr *newrelic.NewRelic,
	article *controller.Article,
	user *controller.User,
	health *controller.Health,
) http.Handler {
	ic := connect.WithInterceptors(interceptor)

	routes := []router.Route{
		router.NewRoute(articlev1connect.NewArticleServiceHandler(article, ic)),
		router.NewRoute(healthv1connect.NewHealthServiceHandler(health, ic)),
		router.NewRoute(userv1connect.NewUserServiceHandler(user, ic)),
	}

	if nr != nil {
		for i, route := range routes {
			routes[i] = router.NewRoute(nr.Handle(route.Path, route.Handler))
		}
	}

	mux := router.New(routes...).Mux()

	mux.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker()))

	return h2c.NewHandler(mux, &http2.Server{})
}

const (
	baseURL   = "/api"
	healthURL = "/health"
)

func NewOpenAPI(
	si openapi.ServerInterface,
	cors openapi.MiddlewareFunc,
	middleware *middleware.Middleware,
) http.Handler {
	router := chi.NewRouter()

	router.Use(cors)

	router.Get(healthURL, func(w http.ResponseWriter, r *http.Request) {})

	return openapi.HandlerWithOptions(si, openapi.ChiServerOptions{
		BaseURL:     baseURL,
		BaseRouter:  router,
		Middlewares: []openapi.MiddlewareFunc{middleware.Handle},
	})
}
