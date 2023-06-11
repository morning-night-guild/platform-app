package notice

type Subject string

func (sub Subject) String() string {
	return string(sub)
}

func NewSubject(value string) (Subject, error) {
	return Subject(value), nil
}

func GenerateInvitationSubject() Subject {
	return Subject("Welcome to Morning Night Guild Platform!")
}
