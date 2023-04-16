package external

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/pkg/trace"
)

const (
	// HeaderTID トレースIDのヘッダー名.
	HeaderTID = "tid"
	// HeaderUID ユーザーIDのヘッダー名.
	HeaderUID = "uid"
)

func NewRequest[T any](ctx context.Context, msg *T) *connect.Request[T] {
	req := connect.NewRequest(msg)

	req.Header().Set(HeaderTID, trace.GetTIDCtx(ctx))

	req.Header().Set(HeaderUID, model.GetUIDCtx(ctx).String())

	return connect.NewRequest(msg)
}
