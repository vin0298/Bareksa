package model

import (
	entity "./entity"
	"github.com/google/uuid"
)

type TagData struct {
	Name string `json:"name"`
}

type NewsArticleData struct {
	Topic         string   `json:"topic"`
	Tags          []string `json:"tags"`
	Title         string   `json:"title"`
	TimePublished string   `json:"timePublished"`
	Content       string   `json:"content"`
	Author        string   `json:"author"`
	Status        string   `json:"status"`
}

type NewsArticle struct {
	topic         entity.Topic
	id            uuid.UUID
	tags          []entity.Tag
	title         string
	timePublished string
	content       string
	author        string
	status        string
}

func NewNewsArticle(articleData NewsArticleData) *NewsArticle {
	newTopic := entity.NewTopic(articleData.Topic)
	tagList := buildTagCollection(articleData.Tags)
	newId := uuid.New()
	return &NewsArticle{
		topic:         newTopic,
		id:            newId,
		tags:          tagList,
		title:         articleData.Title,
		timePublished: articleData.TimePublished,
		content:       articleData.Content,
		author:        articleData.Author,
		status:        articleData.Status,
	}
}

func (n *NewsArticle) Title() string {
	return n.title
}

func (n *NewsArticle) TimePublished() string {
	return n.timePublished
}

func (n *NewsArticle) Content() string {
	return n.content
}

func (n *NewsArticle) Author() string {
	return n.author
}

func (n *NewsArticle) Id() uuid.UUID {
	return n.id
}

func (n *NewsArticle) Tags() []entity.Tag {
	return n.tags
}

func (n *NewsArticle) Status() string {
	return n.status
}

func (n *NewsArticle) GetTopic() string {
	return n.topic.Name()
}

func buildTagCollection(tagList []string) []entity.Tag {
	var tags []entity.Tag
	for _, tagName := range tagList {
		tags = append(tags, entity.NewTag(tagName))
	}
	return tags
}

func CreateArticleTag(tagName string) entity.Tag {
	return entity.NewTag(tagName)
}
