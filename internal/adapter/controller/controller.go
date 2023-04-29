package controller

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var (
	errInternal = fmt.Errorf("internal server")
	ErrInternal = connect.NewError(
		connect.CodeInternal,
		errInternal,
	)
	errInvalidArgument = fmt.Errorf("bad request")
	ErrInvalidArgument = connect.NewError(
		connect.CodeInvalidArgument,
		errInvalidArgument,
	)
	errUnauthorized = fmt.Errorf("unauthorized")
	ErrUnauthorized = connect.NewError(
		connect.CodeUnauthenticated,
		errUnauthorized,
	)
	errNotFound = fmt.Errorf("not found")
	ErrNotFound = connect.NewError(
		connect.CodeNotFound,
		errNotFound,
	)
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

// HandleConnectError 発生したエラーを対応するconnectのcode込みのエラーに変換する関数.
func (ctrl *Controller) HandleConnectError(ctx context.Context, err error) error {
	logger := log.GetLogCtx(ctx)

	switch {
	case
		errors.AsValidationError(err),
		errors.AsURLError(err):
		logger.Warn(err.Error())

		return ErrInvalidArgument
	case errors.AsNotFoundError(err):
		logger.Warn(err.Error())

		return ErrNotFound
	default:
		logger.Error(err.Error())

		return ErrInternal
	}
}
