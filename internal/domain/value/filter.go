package value

type Filter struct {
	Name  string
	Value string // 現状文字列のみ対応 ジェネリクス対応したかったが関数の引数に渡し方が不明なので一旦文字列
}

func NewFilter(
	name string,
	value string,
) Filter {
	return Filter{
		Name:  name,
		Value: value,
	}
}
