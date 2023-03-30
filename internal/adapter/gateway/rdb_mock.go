package gateway

import (
	"fmt"
	"testing"

	"github.com/morning-night-guild/platform-app/pkg/ent"
	"github.com/morning-night-guild/platform-app/pkg/ent/enttest"
)

var _ RDBFactory = (*RDBClientMock)(nil)

type RDBClientMock struct {
	t *testing.T
}

func NewRDBClientMock(t *testing.T) *RDBClientMock {
	t.Helper()

	return &RDBClientMock{
		t: t,
	}
}

func (rm *RDBClientMock) Of(dsn string) (*RDB, error) {
	rm.t.Helper()

	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(rm.t.Log)),
	}

	dataSourceName := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", dsn)

	db := enttest.Open(rm.t, "sqlite3", dataSourceName, opts...)

	return &RDB{
		Client: db,
	}, nil
}
