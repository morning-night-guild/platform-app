package rpc

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/notice"
)

//go:generate mockgen -source notice.go -destination notice_mock.go -package rpc

type Notice interface {
	Notify(context.Context, auth.Email, notice.Subject, notice.Message) (notice.ID, error)
}
