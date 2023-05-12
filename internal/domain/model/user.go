package model

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

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
