package router

import (
	"net/http"
)

// Router.
type Router struct {
	routes []Route
}

// NewRouter.
func New(routes ...Route) *Router {
	return &Router{
		routes: routes,
	}
}

// Mux.
func (r Router) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range r.routes {
		mux.Handle(route.Path, route.Handler)
	}

	return mux
}

// Route.
type Route struct {
	Path    string
	Handler http.Handler
}

// NewRoute.
func NewRoute(path string, handler http.Handler) Route {
	return Route{
		Path:    path,
		Handler: handler,
	}
}
