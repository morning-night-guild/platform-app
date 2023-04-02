package gateway_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/adapter/gateway"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func TestUserSave(t *testing.T) {
	t.Parallel()

	t.Run("ユーザーの保存できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		userGateway := gateway.NewUser(rdb)

		ctx := context.Background()

		usr := model.CreateUser()

		if err := userGateway.Save(ctx, usr); err != nil {
			t.Error(err)
		}

		found, err := userGateway.Find(ctx, usr.UserID)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(found, usr) {
			t.Errorf("User() = %v, want %v", found, usr)
		}
	})

	t.Run("同じユーザーを複数回保存できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		userGateway := gateway.NewUser(rdb)

		ctx := context.Background()

		usr := model.CreateUser()

		if err := userGateway.Save(ctx, usr); err != nil {
			t.Error(err)
		}

		found1, err := rdb.User.Get(ctx, usr.UserID.Value())
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)

		if err := userGateway.Save(ctx, usr); err != nil {
			t.Error(err)
		}

		found2, err := rdb.User.Get(ctx, usr.UserID.Value())
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(found1.ID, found2.ID) {
			t.Errorf("User() id = %v, want %v", found1, found2)
		}

		if !reflect.DeepEqual(found1.CreatedAt, found2.CreatedAt) {
			t.Errorf("User() crated_at = %v, want %v", found1.CreatedAt, found2.CreatedAt)
		}

		// 更新日時がfound1よりfound2の方があとであることを確認する
		if found1.UpdatedAt.After(found2.UpdatedAt) {
			t.Errorf("User() updated_at found1 = %v, found2 = %v", found1.UpdatedAt, found2.UpdatedAt)
		}
	})
}

func TestUserFind(t *testing.T) {
	t.Parallel()

	t.Run("ユーザーが取得できる", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		userGateway := gateway.NewUser(rdb)

		ctx := context.Background()

		usr := model.CreateUser()

		if err := userGateway.Save(ctx, usr); err != nil {
			t.Error(err)
		}

		found, err := userGateway.Find(ctx, usr.UserID)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(found, usr) {
			t.Errorf("User() = %v, want %v", found, usr)
		}
	})

	t.Run("ユーザーが取得できない", func(t *testing.T) {
		t.Parallel()

		rdb, err := gateway.NewRDBClientMock(t).Of(uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}

		userGateway := gateway.NewUser(rdb)

		ctx := context.Background()

		uid := user.GenerateID()

		if _, err := userGateway.Find(ctx, uid); err == nil {
			t.Error("found user")
		}
	})
}
