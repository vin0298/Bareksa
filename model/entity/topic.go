package entity

type Topic struct {
	name string
}

func NewTopic(topicName string) {
	return Topic{name: topicName}
}
