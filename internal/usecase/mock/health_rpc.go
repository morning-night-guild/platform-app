package mock

import (
	"context"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
)

var _ rpc.Health = (*RPCHealth)(nil)

type RPCHealth struct {
	T        *testing.T
	CheckErr error
}

func (rh *RPCHealth) Check(ctx context.Context) error {
	rh.T.Helper()

	return rh.CheckErr
}
