package notice

type ID string

func (id ID) String() string {
	return string(id)
}

func NewID(id string) (ID, error) {
	return ID(id), nil
}
