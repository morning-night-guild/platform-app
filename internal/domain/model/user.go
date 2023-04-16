package model

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

type uky struct{}

type User struct {
	UserID user.ID
}

func NewUser(
	userID user.ID,
) User {
	return User{
		UserID: userID,
	}
}

func CreateUser() User {
	id := user.GenerateID()

	return NewUser(id)
}

func SetUIDCtx(ctx context.Context, uid user.ID) context.Context {
	return context.WithValue(ctx, uky{}, uid)
}

func GetUIDCtx(ctx context.Context) user.ID {
	v := ctx.Value(uky{})

	uid, ok := v.(user.ID)
	if !ok {
		return user.GenerateZeroID()
	}

	return uid
}
