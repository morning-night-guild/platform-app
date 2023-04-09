package user

import "github.com/google/uuid"

type ID uuid.UUID

// NewID IDを作成するファクトリー関数.
func NewID(value string) (ID, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return ID{}, err
	}

	return ID(uid), nil
}

// GenerateID IDを新規に発行する関数.
func GenerateID() ID {
	return ID(uuid.New())
}

func GenerateZeroID() ID {
	return ID(uuid.Nil)
}

// Value IDをuuid.UUID型として提供するメソッド.
func (uid ID) Value() uuid.UUID {
	return uuid.UUID(uid)
}

// String IDを文字列型として提供するメソッド.
func (uid ID) String() string {
	return uid.Value().String()
}
