package entity

type Tag struct {
	name string
}

func NewTag(tagName string) {
	return Tag{name: tagName}
}
