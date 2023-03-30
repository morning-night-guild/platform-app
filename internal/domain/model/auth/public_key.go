package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

func DecodePublicKey(value string) (rsa.PublicKey, error) {
	pubkey, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		msg := "failed to decode public key"

		return rsa.PublicKey{}, errors.NewValidationError(msg)
	}

	pub, err := x509.ParsePKIXPublicKey(pubkey)
	if err != nil {
		msg := "failed to parse public key"

		return rsa.PublicKey{}, errors.NewValidationError(msg)
	}

	key, ok := pub.(*rsa.PublicKey)
	if !ok {
		msg := "failed to convert public key"

		return rsa.PublicKey{}, errors.NewValidationError(msg)
	}

	return *key, nil
}
