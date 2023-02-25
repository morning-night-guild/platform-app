package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

var _ repository.APIHealth = (*APIHealth)(nil)

type APIHealth struct {
	T   *testing.T
	Err error
}

func (ah *APIHealth) Check(ctx context.Context) error {
	ah.T.Helper()

	return ah.Err
}
