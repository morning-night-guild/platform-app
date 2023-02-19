package cors

import (
	"errors"
	"strings"

	"github.com/go-chi/cors"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

var (
	errEmptyArray  = errors.New("empty array")
	errEmptyString = errors.New("empty string")
)

const maxAge = 300

func New(allowOrigins []string, debug bool) (openapi.MiddlewareFunc, error) {
	if len(allowOrigins) == 0 {
		return nil, errEmptyArray
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   allowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		Debug:            debug,
		MaxAge:           maxAge,
	}), nil
}

func ConvertAllowOrigins(allowOrigins string) ([]string, error) {
	if allowOrigins == "" {
		return nil, errEmptyString
	}

	return strings.Split(allowOrigins, ","), nil
}

func ConvertDebugEnable(debug string) bool {
	return debug == "true"
}
