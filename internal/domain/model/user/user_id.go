package user

import "github.com/google/uuid"

type UserID uuid.UUID //nolint:revive

// NewUserID UserIDを作成するファクトリー関数.
func NewUserID(value string) (UserID, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, err
	}

	return UserID(uid), nil
}

// GenerateUserID UserIDを新規に発行する関数.
func GenerateUserID() UserID {
	return UserID(uuid.New())
}

// Value UserIDをuuid.UUID型として提供するメソッド.
func (uid UserID) Value() uuid.UUID {
	return uuid.UUID(uid)
}

// String UserIDを文字列型として提供するメソッド.
func (uid UserID) String() string {
	return uid.Value().String()
}
