package article

// Title 記事のタイトル.
type Title string

// String 記事のタイトルを文字列として提供するメソッド.
func (t Title) String() string {
	return string(t)
}

// NewTitle 記事のタイトルを新規作成する関数.
func NewTitle(t string) (Title, error) {
	title := Title(t)

	if err := title.validate(); err != nil {
		return Title(""), err
	}

	return title, nil
}

// validate 記事のタイトルを検証するメソッド.
func (t Title) validate() error {
	return nil
}
