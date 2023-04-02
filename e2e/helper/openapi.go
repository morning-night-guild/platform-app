package helper

import (
	"net/http"
	"testing"

	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

type OpenAPIClient struct {
	Client *openapi.Client
}

func NewOpenAPIClient(t *testing.T, url string) *OpenAPIClient {
	t.Helper()

	client, err := openapi.NewClient(url + "/api")
	if err != nil {
		t.Error(err)
	}

	return &OpenAPIClient{
		Client: client,
	}
}

func NewOpenAPIClientWithAPIKey(t *testing.T, url string, key string) *OpenAPIClient {
	t.Helper()

	client, err := openapi.NewClient(url + "/api")
	if err != nil {
		t.Error(err)
	}

	client.Client = &http.Client{
		Transport: NewAPIKeyTransport(t, key),
	}

	return &OpenAPIClient{
		Client: client,
	}
}
