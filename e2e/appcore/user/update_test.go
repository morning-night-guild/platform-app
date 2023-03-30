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

func TestAppCoreE2EUserUpdate(t *testing.T) {
	t.Parallel()

	url := helper.GetAppCoreEndpoint(t)

	t.Run("ユーザーが更新できる", func(t *testing.T) {
		t.Parallel()

		db := helper.NewDatabase(t, helper.GetDSN(t))

		defer db.Close()

		client := helper.NewConnectClient(t, http.DefaultClient, url)

		req1 := &userv1.CreateRequest{}

		res1, err := client.User.Create(context.Background(), connect.NewRequest(req1))
		if err != nil {
			t.Errorf("failed to health check: %s", err)
		}

		req2 := &userv1.UpdateRequest{
			UserId: res1.Msg.User.UserId,
		}

		res2, err := client.User.Update(context.Background(), connect.NewRequest(req2))
		if err != nil {
			t.Errorf("failed to health check: %s", err)
		}

		db.DeleteUser(uuid.MustParse(res2.Msg.User.UserId))
	})
}
