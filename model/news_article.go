package model

import (
	entity "./entity"
)

type NewsArticle struct {
	topic          entity.Topic
	tag            entity.Tag
	ID             string
	title          string
	timePublished  string
	content        string
	author         string
	newsRepository NewsArticleRepository
}
