package handler

import "net/http"

//go:generate mockgen -source cookie.go -destination cookie_mock.go -package handler

type Cookie interface {
	Domain() string
	Secure() bool
	SameSite() http.SameSite
}
