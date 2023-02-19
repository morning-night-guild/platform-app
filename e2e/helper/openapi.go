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

type APIKeyTransport struct {
	t         *testing.T
	APIKey    string
	Transport http.RoundTripper
}

func NewAPIKeyTransport(
	t *testing.T,
	key string,
) *APIKeyTransport {
	t.Helper()

	return &APIKeyTransport{
		t:         t,
		APIKey:    key,
		Transport: http.DefaultTransport,
	}
}

func (at *APIKeyTransport) transport() http.RoundTripper {
	at.t.Helper()

	return at.Transport
}

func (at *APIKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	at.t.Helper()

	req.Header.Add("Api-Key", at.APIKey)

	resp, err := at.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
