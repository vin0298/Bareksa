package model

import (
	entity "./entity"
	"github.com/google/uuid"
)

type NewsArticle struct {
	topic          entity.Topic
	ID			   UUID
	tags            []string
	title          string
	timePublished  string
	content        string
	author         string
}


func NewNewsArticle(articleData map[string]interface{}) NewsArticle {
	newTopic := entity.NewTopic(articleData["topic"])
	return &NewsArticle {
		topic: newTopic,
		ID: uuid.New(),
		tags: articleData["tags"],
		title: articleData["title"],
		timePublished: articleData["time"],
		content: articleData["content"],
		author: articleData["author"]
	}
}
