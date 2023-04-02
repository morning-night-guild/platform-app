package auth

import "github.com/morning-night-guild/platform-app/internal/domain/model/errors"

type Signature string

func NewSignature(value string) (Signature, error) {
	sig := Signature(value)

	if err := sig.validate(); err != nil {
		return Signature(""), err
	}

	return sig, nil
}

func (sig Signature) String() string {
	return string(sig)
}

func (sig Signature) validate() error {
	if sig == "" {
		return errors.NewValidationError("signature is required")
	}

	return nil
}
