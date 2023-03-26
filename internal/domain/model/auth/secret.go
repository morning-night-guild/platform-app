package auth

type Secret string

func NewSecret(secret string) Secret {
	return Secret(secret)
}

func (s Secret) String() string {
	return string(s)
}
