package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/e2e/helper"
	userv1 "github.com/morning-night-guild/platform-app/pkg/connect/user/v1"
)

func TestAppCoreE2EUserCreate(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("ユーザーが作成できる", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		client := helper.NewConnectClient(t, http.DefaultClient, url)

		req := &userv1.CreateRequest{}

		res, err := client.User.Create(context.Background(), connect.NewRequest(req))
		if err != nil {
			t.Errorf("failed to health check: %s", err)
		}

		db.DeleteUser(uuid.MustParse(res.Msg.User.UserId))
	})
}
