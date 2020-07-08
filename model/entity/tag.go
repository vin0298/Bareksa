package entity

import "github.com/google/uuid"

type Tag struct {
	name string
	id   uuid.UUID
}

func NewTag(tagName string) Tag {
	return Tag{name: tagName, id: uuid.New()}
}

func (t Tag) Name() string {
	return t.name
}

func (t Tag) Id() uuid.UUID {
	return t.id
}
