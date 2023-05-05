package external

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
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

	return req
}

type External struct{}

func New() *External {
	return &External{}
}

func (ext *External) HandleError(_ context.Context, err error) error {
	code := connect.CodeOf(err)
	switch code { //nolint:exhaustive
	case connect.CodeInvalidArgument:
		return errors.NewValidationError(code.String(), err)
	case connect.CodeNotFound:
		return errors.NewNotFoundError(code.String(), err)
	case connect.CodeUnauthenticated:
		return errors.NewUnauthorizedError(code.String(), err)
	default:
		return errors.NewUnknownError(code.String(), err)
	}
}
