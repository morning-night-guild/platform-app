package auth

import "github.com/google/uuid"

type SessionID uuid.UUID

// NewSessionID SessionIDを作成するファクトリー関数.
func NewSessionID(value string) (SessionID, error) {
	sid, err := uuid.Parse(value)
	if err != nil {
		return SessionID{}, err
	}

	return SessionID(sid), nil
}

// GenerateSessionID SessionIDを新規に発行する関数.
func GenerateSessionID() SessionID {
	return SessionID(uuid.New())
}

// Value SessionIDをuuid.UUID型として提供するメソッド.
func (sid SessionID) Value() uuid.UUID {
	return uuid.UUID(sid)
}

// String SessionIDを文字列型として提供するメソッド.
func (sid SessionID) String() string {
	return sid.Value().String()
}

// ToSecret SessionIDをSecret型として提供するメソッド.
func (sid SessionID) ToSecret() Secret {
	return Secret(sid.String())
}
