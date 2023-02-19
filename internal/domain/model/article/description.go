package article

// Description 記事の説明.
type Description string

// NewDescription 記事の説明を新規作成する関数.
func NewDescription(value string) (Description, error) {
	description := Description(value)

	if err := description.validate(); err != nil {
		return Description(""), err
	}

	return description, nil
}

// String 記事の説明を文字列として提供するメソッド.
func (d Description) String() string {
	return string(d)
}

// validate 記事の説明を検証するメソッド.
func (d Description) validate() error {
	return nil
}
