package article

import (
	"fmt"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

// Tag タグリスト.
type TagList []Tag

// maxTagLength タグリストに含まれるタグの個数の最大値.
const maxTagLength = 5

// NewTagList タグリストを新規作成する関数.
func NewTagList(values []Tag) (TagList, error) {
	tags := TagList(values).distinct()

	if err := tags.validate(); err != nil {
		return nil, err
	}

	return tags, nil
}

// distinct 重複を排除するメソッド.
func (t TagList) distinct() TagList {
	tmp := make(map[Tag]struct{}, t.Len())

	for _, v := range t {
		tmp[v] = struct{}{}
	}

	uniq := make([]Tag, 0, t.Len())
	for t := range tmp {
		uniq = append(uniq, t)
	}

	return TagList(uniq)
}

// Contains リスト内のタグの存在チェックを行うメソッド.
func (t TagList) Contains(target Tag) bool {
	for _, tag := range t {
		if tag.String() == target.String() {
			return true
		}
	}

	return false
}

// Append タグを追加するメソッド.
func (t TagList) Append(tag Tag) TagList {
	if t.Contains(tag) {
		return t
	}

	return append(t, tag)
}

// validate タグリストを検証するメソッド.
func (t TagList) validate() error {
	if t.Len() > maxTagLength {
		msg := fmt.Sprintf("must be less than or equal to %d. length is %d", maxTagLength, t.Len())

		return errors.NewValidationError(msg)
	}

	return nil
}

// Len タグリストに含まれるのタグの個数を提供するメソッド.
func (t TagList) Len() int {
	return len(t)
}

// StringSlice 文字列型のスライスを提供するメソッド.
func (t TagList) StringSlice() []string {
	list := make([]string, 0, t.Len())

	for _, tag := range t {
		list = append(list, tag.String())
	}

	return list
}
