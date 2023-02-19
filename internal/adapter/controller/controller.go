package controller

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	me "github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var errInternal = errors.New("internal server")

var ErrInternal = connect.NewError(
	connect.CodeInternal,
	errInternal,
)

var errInvalidArgument = errors.New("bad request")

var ErrInvalidArgument = connect.NewError(
	connect.CodeInvalidArgument,
	errInvalidArgument,
)
var errUnauthorized = errors.New("unauthorized")

var ErrUnauthorized = connect.NewError(
	connect.CodeUnauthenticated,
	errUnauthorized,
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

// HandleConnectError 発生したエラーを対応するconnectのcode込みのエラーに変換する関数.
func (ctl *Controller) HandleConnectError(ctx context.Context, err error) error {
	logger := log.GetLogCtx(ctx)

	switch {
	case
		ctl.asValidationError(err),
		ctl.asURLError(err):
		logger.Warn(err.Error())

		return ErrInvalidArgument
	default:
		logger.Error(err.Error())

		return ErrInternal
	}
}

func (ctl *Controller) asValidationError(err error) bool {
	var target me.ValidationError

	return errors.As(err, &target)
}

func (ctl *Controller) asURLError(err error) bool {
	var target me.URLError

	return errors.As(err, &target)
}
