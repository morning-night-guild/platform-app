package resend

import (
	"github.com/morning-night-guild/platform-app/internal/adapter/external"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/resendlabs/resend-go"
)

var _ external.NoticeFactory = (*Resend)(nil)

type Resend struct{}

func New() *Resend {
	return &Resend{}
}

func (rsd *Resend) Notice( //nolint:ireturn
	key string,
	from string,
) (rpc.Notice, error) {
	cl := resend.NewClient(key)

	return external.NewNotice(
		cl,
		from,
	), nil
}

func (rsd *Resend) MockNotice() rpc.Notice { //nolint:ireturn
	return external.NewMockNotice()
}
