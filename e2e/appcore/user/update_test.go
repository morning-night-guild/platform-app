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
			t.Errorf("failed to create user: %s", err)
		}

		req2 := &userv1.UpdateRequest{
			UserId: res1.Msg.User.UserId,
		}

		res2, err := client.User.Update(context.Background(), connect.NewRequest(req2))
		if err != nil {
			t.Errorf("failed to update user: %s", err)
		}

		db.DeleteUser(uuid.MustParse(res2.Msg.User.UserId))
	})

	t.Run("存在しないユーザーであるため更新に失敗する", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, http.DefaultClient, url)

		req := &userv1.UpdateRequest{
			UserId: uuid.New().String(),
		}

		if _, err := client.User.Update(context.Background(), connect.NewRequest(req)); err == nil {
			t.Errorf("success to update user: %s", err)
		}
	})

	t.Run("不正なIDであるため更新に失敗する", func(t *testing.T) {
		t.Parallel()

		client := helper.NewConnectClient(t, http.DefaultClient, url)

		req := &userv1.UpdateRequest{
			UserId: "uid",
		}

		if _, err := client.User.Update(context.Background(), connect.NewRequest(req)); err == nil {
			t.Errorf("success to update user: %s", err)
		}
	})
}
