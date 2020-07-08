package entity

type Topic struct {
	name string
}

func NewTopic(topicName string) Topic {
	return Topic{name: topicName}
}

func(t Topic) Name() string {
	return t.name
}
