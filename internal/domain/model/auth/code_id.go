package auth

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

type CodeID uuid.UUID

// NewCodeID CodeIDを作成するファクトリー関数.
func NewCodeID(value string) (CodeID, error) {
	cd, err := uuid.Parse(value)
	if err != nil {
		return CodeID{}, err
	}

	return CodeID(cd), nil
}

// GenerateCodeID CodeIDを新規に発行する関数.
func GenerateCodeID() CodeID {
	return CodeID(uuid.New())
}

// Value CodeIDをuuid.UUID型として提供するメソッド.
func (cd CodeID) Value() uuid.UUID {
	return uuid.UUID(cd)
}

// String CodeIDを文字列型として提供するメソッド.
func (cd CodeID) String() string {
	return cd.Value().String()
}

func (cd CodeID) Verify(
	signature Signature,
	publicKey *rsa.PublicKey,
) error {
	h := crypto.Hash.New(crypto.SHA256)

	h.Write([]byte(cd.String()))

	hashed := h.Sum(nil)

	sig, err := base64.StdEncoding.DecodeString(signature.String())
	if err != nil {
		return errors.NewValidationError("Signature is invalid")
	}

	if err := rsa.VerifyPSS(publicKey, crypto.SHA256, hashed, sig, &rsa.PSSOptions{
		Hash: crypto.SHA256,
	}); err != nil {
		return errors.NewValidationError("Signature is invalid")
	}

	return nil
}
