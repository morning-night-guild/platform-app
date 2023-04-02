package auth

type Secret string

func (s Secret) String() string {
	return string(s)
}
