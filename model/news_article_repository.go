package model

type NewsArticleRepository interface {
	FindNewsByTopic() ([]NewsArticle, error)
	FindNewsByStatus() ([]NewsArticle, error)
	RetrieveAnArticle(articleID string) (NewsArticle, error)
	DeleteAnArticle(articleID string) error
	CreateAnArticle(NewsArticle) error
	UpdateAnArticle(NewsArticle) error
}
