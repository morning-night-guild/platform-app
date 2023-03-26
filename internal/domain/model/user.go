package model

import "github.com/morning-night-guild/platform-app/internal/domain/model/user"

type User struct {
	UserID user.UserID
}

func NewUser(
	userID user.UserID,
) User {
	return User{
		UserID: userID,
	}
}

func CreateUser() User {
	id := user.GenerateUserID()

	return NewUser(id)
}
