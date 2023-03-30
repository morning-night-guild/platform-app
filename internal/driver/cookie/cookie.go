package cookie

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
)

var _ handler.Cookie = (*Cookie)(nil)

type Cookie struct{}

func New() *Cookie {
	return &Cookie{}
}

func (ck *Cookie) Domain() string {
	return "localhost"
}

func (ck *Cookie) HTTPOnly() bool {
	return true
}

func (ck *Cookie) SameSite() http.SameSite {
	return http.SameSiteDefaultMode
}

func (ck *Cookie) Secure() bool {
	return true
}
