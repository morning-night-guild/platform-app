package newrelic

import (
	"fmt"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewRelic.
type NewRelic struct {
	newrelic *newrelic.Application
}

// New.
func New(
	appName string,
	license string,
) (*NewRelic, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new relic application: %w", err)
	}

	return &NewRelic{
		newrelic: app,
	}, nil
}

// Handle.
// NewRelicへアクセスログを送信するためのミドルウェアハンドラ.
func (nr *NewRelic) Handle(
	pattern string,
	handler http.Handler,
	options ...newrelic.TraceOption,
) (string, http.Handler) {
	return newrelic.WrapHandle(nr.newrelic, pattern, handler, options...)
}
