package firebase

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/morning-night-guild/platform-app/internal/adapter/external"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"google.golang.org/api/option"
)

type Firebase struct{}

func New() *Firebase {
	return &Firebase{}
}

func (fb *Firebase) Of(
	secret string,
	endpoint string,
	apiKey string,
) (*external.Auth, error) {
	opt := option.WithCredentialsJSON([]byte(secret))

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Log().Error("error initializing app", log.ErrorField(err))

		return nil, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Log().Error("error create auth client", log.ErrorField(err))

		return nil, err
	}

	if endpoint == "" {
		return nil, errors.NewValidationError("firebase endpoint is empty")
	}

	if apiKey == "" {
		return nil, errors.NewValidationError("firebase api key is empty")
	}

	return external.NewAuth(
		endpoint,
		apiKey,
		http.DefaultClient,
		auth,
	), nil
}
