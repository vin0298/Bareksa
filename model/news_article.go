package model

import (
	entity "./entity"
	"github.com/google/uuid"
)

type NewsArticleData struct {
	Topic			string `json:"topic"`
	Tags 			[]string `json:"tags"`
	Title 			string `json:"title"`
	TimePublished 	string `json:"timePublished`
	Content 		string `json:"content"`
	Author  		string `json:"author"`
}

type NewsArticle struct {
	topic          entity.Topic
	id			   uuid.UUID
	tags            []string
	title          string
	timePublished  string
	content        string
	author         string
}

func NewNewsArticle(articleData NewsArticleData) *NewsArticle {
	newTopic := entity.NewTopic(articleData.Topic)
	newId := uuid.New()
	return &NewsArticle {
		topic: newTopic,
		id: newId,
		tags: articleData.Tags,
		title: articleData.Title,
		timePublished: articleData.TimePublished,
		content: articleData.Content,
		author: articleData.Author,
	}
}

func (n * NewsArticle) Title() string {
	return n.title
}

func (n * NewsArticle) TimePublished() string {
	return n.time
}

func (n * NewsArticle) Content() string {
	return n.content
}

func (n * NewsArticle) Author() string {
	return n.author
}

func (n * NewsArticle) ID
