package model

type ArticleReadModel struct {
	Author         string
	Title          string
	Content        string
	Time_published string
	Uuid           string
	Topic          string
	Status         string
	Tags           []string
}
