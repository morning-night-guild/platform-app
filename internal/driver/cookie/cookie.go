package cookie

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/adapter/handler"
)

var _ handler.Cookie = (*Cookie)(nil)

type Cookie struct {
	domain string
}

func New(
	domain string,
) *Cookie {
	return &Cookie{
		domain: domain,
	}
}

func (ck *Cookie) Domain() string {
	return ck.domain
}

func (ck *Cookie) SameSite() http.SameSite {
	if ck.domain == "" {
		return http.SameSiteDefaultMode
	}

	return http.SameSiteNoneMode
}

func (ck *Cookie) Secure() bool {
	return ck.domain != ""
}
