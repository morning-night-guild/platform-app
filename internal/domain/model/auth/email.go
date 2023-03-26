package auth

type EMail string

func NewEMail(value string) (EMail, error) {
	em := EMail(value)

	if err := em.validate(); err != nil {
		return EMail(""), err
	}

	return em, nil
}

func (em EMail) String() string {
	return string(em)
}

func (em EMail) validate() error {
	return nil
}
