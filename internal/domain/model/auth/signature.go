package auth

type Signature string

func NewSignature(value string) (Signature, error) {
	return Signature(value), nil
}

func (sig Signature) String() string {
	return string(sig)
}
