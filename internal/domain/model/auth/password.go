package auth

type Password string

func NewPassword(value string) (Password, error) {
	pw := Password(value)

	if err := pw.validate(); err != nil {
		return Password(""), err
	}

	return pw, nil
}

func (pw Password) String() string {
	return string(pw)
}

func (pw Password) validate() error {
	return nil
}
