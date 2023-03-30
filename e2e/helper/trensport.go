package helper

import (
	"net/http"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

type CookieTransport struct {
	t         *testing.T
	Cookie    string
	Transport http.RoundTripper
}

func NewCookieTransport(
	t *testing.T,
	cookie string,
) *CookieTransport {
	return &CookieTransport{
		t:         t,
		Cookie:    cookie,
		Transport: http.DefaultTransport,
	}
}

func (ct *CookieTransport) transport() http.RoundTripper {
	ct.t.Helper()

	return ct.Transport
}

func (ct *CookieTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ct.t.Helper()

	req.Header.Add("Cookie", ct.Cookie)

	resp, err := ct.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

type CookiesTransport struct {
	T         *testing.T
	Cookies   []*http.Cookie
	Transport http.RoundTripper
}

func NewCookiesTransport(
	t *testing.T,
	cookies []*http.Cookie,
) *CookiesTransport {
	t.Helper()

	return &CookiesTransport{
		T:         t,
		Cookies:   cookies,
		Transport: http.DefaultTransport,
	}
}

func (ct *CookiesTransport) transport() http.RoundTripper {
	ct.T.Helper()

	return ct.Transport
}

func (ct *CookiesTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ct.T.Helper()

	for _, cookie := range ct.Cookies {
		req.AddCookie(cookie)
	}

	resp, err := ct.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

type OnlyIDTokenCookieTransport struct {
	T         *testing.T
	Cookie    *http.Cookie
	Transport http.RoundTripper
}

func NewOnlyIDTokenCookieTransport(
	t *testing.T,
	cookies []*http.Cookie,
) *OnlyIDTokenCookieTransport {
	t.Helper()

	var cookie *http.Cookie

	for _, c := range cookies {
		if c.Name == auth.AuthTokenKey {
			cookie = c
			break
		}
	}

	return &OnlyIDTokenCookieTransport{
		T:         t,
		Cookie:    cookie,
		Transport: http.DefaultTransport,
	}
}

func (ct *OnlyIDTokenCookieTransport) transport() http.RoundTripper {
	ct.T.Helper()

	return ct.Transport
}

func (ct *OnlyIDTokenCookieTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ct.T.Helper()

	req.AddCookie(ct.Cookie)

	resp, err := ct.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

type OnlySessionTokenCookieTransport struct {
	T         *testing.T
	Cookie    *http.Cookie
	Transport http.RoundTripper
}

func NewOnlySessionTokenCookieTransport(
	t *testing.T,
	cookies []*http.Cookie,
) *OnlySessionTokenCookieTransport {
	t.Helper()

	var cookie *http.Cookie

	for _, c := range cookies {
		if c.Name == auth.SessionTokenKey {
			cookie = c
			break
		}
	}

	return &OnlySessionTokenCookieTransport{
		T:         t,
		Cookie:    cookie,
		Transport: http.DefaultTransport,
	}
}

func (ct *OnlySessionTokenCookieTransport) transport() http.RoundTripper {
	ct.T.Helper()

	return ct.Transport
}

func (ct *OnlySessionTokenCookieTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ct.T.Helper()

	req.AddCookie(ct.Cookie)

	resp, err := ct.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, err
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
