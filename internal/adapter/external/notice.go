package external

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/notice"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/resendlabs/resend-go"
)

type NoticeFactory interface {
	Notice(key, from string) (rpc.Notice, error)
	MockNotice() rpc.Notice
}

var (
	_ rpc.Notice = (*Notice)(nil)
	_ rpc.Notice = (*MockNotice)(nil)
)

type Notice struct {
	client *resend.Client
	from   string
}

func NewNotice( //nolint:ireturn
	client *resend.Client,
	from string,
) rpc.Notice {
	return &Notice{
		client: client,
		from:   from,
	}
}

func (ntc *Notice) Notify(
	ctx context.Context,
	to auth.Email,
	subject notice.Subject,
	message notice.Message,
) (notice.ID, error) {
	req := &resend.SendEmailRequest{
		From:    ntc.from,
		To:      []string{to.String()},
		Subject: subject.String(),
		Text:    message.String(),
	}

	res, err := ntc.client.Emails.Send(req)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to send email", log.ErrorField(err))

		return notice.ID(""), err
	}

	return notice.ID(res.Id), nil
}

type MockNotice struct{}

func NewMockNotice() rpc.Notice { //nolint:ireturn
	return &MockNotice{}
}

func (lg *MockNotice) Notify(
	ctx context.Context,
	to auth.Email,
	subject notice.Subject,
	message notice.Message,
) (notice.ID, error) {
	format := "send email to %s, subject: %s, message: %s"

	log.GetLogCtx(ctx).Info(fmt.Sprintf(format, to.String(), subject.String(), message.String()))

	return notice.ID(uuid.NewString()), nil
}
