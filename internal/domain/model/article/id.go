package article

import "github.com/google/uuid"

// ID.
type ID uuid.UUID

// NewID IDを作成するファクトリー関数.
func NewID(value string) (ID, error) {
	i, err := uuid.Parse(value)
	if err != nil {
		return ID{}, err
	}

	return ID(i), nil
}

// GenerateID IDを新規に発行する関数.
func GenerateID() ID {
	return ID(uuid.New())
}

// Value IDをuuid.UUID型として提供するメソッド.
func (i ID) Value() uuid.UUID {
	return uuid.UUID(i)
}

// String IDを文字列型として提供するメソッド.
func (i ID) String() string {
	return i.Value().String()
}
