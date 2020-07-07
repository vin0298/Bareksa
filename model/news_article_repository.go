package model

import (
	"log"

	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"database/sql"
	_ "github.com/lib/pq"
)

type NewsArticleRepository interface {
	FindNewsByTopic() ([]NewsArticle, error)
	FindNewsByStatus() ([]NewsArticle, error)
	RetrieveAnArticle(articleID string) (*NewsArticle, error)
	DeleteAnArticle(articleID string) error
	CreateAnArticle(*NewsArticle) error
	UpdateAnArticle(*NewsArticle) error
}

type newsArticleRepository struct {
	db               *sql.DB
}

/* TEMP */
func SetupDatabase() (NewsArticleRepository, error) {
	/* Enable SSL later */
	psqlInfo :=  fmt.Sprintf("host=%s port=%d user=%s "+
					"password=%s dbname=%s sslmode=disable",
					viper.GetString("DB_HOST"), viper.GetInt("DB_PORT"), 
					viper.GetString("DB_USER"), viper.GetString("DB_PW"), 
					viper.GetString("DB_NAME"))

	postgresDB, err := sql.Open("postgres", psqlInfo)
	if (err != nil) {
		log.Printf("%s", err)
		return newsArticleRepository{}, err
	}

	err = postgresDB.Ping()
	if (err != nil) {
		log.Printf("%s", err)
		return newsArticleRepository{}, err
	}

	return newsArticleRepository{db: postgresDB}, nil
}

func (n newsArticleRepository) FindNewsByTopic() ([]NewsArticle, error) {
	return nil, nil
}

func (n newsArticleRepository) FindNewsByStatus() ([]NewsArticle, error) {
	return nil, nil
}

func (n newsArticleRepository) RetrieveAnArticle(articleID string) (*NewsArticle, error) {
	return &NewsArticle{}, nil
}

func (n newsArticleRepository) DeleteAnArticle(articleID string) error {
	return nil
}

func (n newsArticleRepository) CreateAnArticle(article *NewsArticle) error {
	sqlStatement := `INSERT INTO news_articles(author, title, content, time_published, uuid, topic)
					 VALUES ($1, $2, $3, $4, $5, $6) RETURNING article_id`
	
	/* Insert the NewsArticle content */
	articleId := 0
	err := n.db.QueryRow(sqlStatement, article.Author(), article.Title(), 
				article.Content(), article.TimePublished(), article.Id(), 
				article.GetTopic()).Scan(&articleId)		
	
	if err != nil {
		log.Printf("Error when inserting a new NewsArticle: %s", err)
		return err
	}

	/* Refactor: Insert the tags */
	sqlStatement = `INSERT INTO tags(tag_name, uuid) VALUES`
	tagValues := []interface{}{}
	for i, tag := range article.Tags() {
		tagValues = append(tagValues, tag.Name(), tag.Id())
		numFields := 2

		n := i * numFields
		sqlStatement += `(`
		for j := 0; j < numFields; j++ {
			sqlStatement += `$` + strconv.Itoa(n + j + 1) + `,`
		}
		sqlStatement = sqlStatement[:len(sqlStatement) - 1] + `),`
	}
	sqlStatement = sqlStatement[:len(sqlStatement) - 1]
	sqlStatement += `ON CONFLICT DO NOTHING RETURNING tag_id`

	rows, err := n.db.Query(sqlStatement, tagValues...)
	defer rows.Close()

	var tagIdList []int

	for rows.Next() {
		var tag_id int 
		err = rows.Scan(&tag_id)
		if err != nil {
			return err
		}
		tagIdList = append(tagIdList, tag_id)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	/* End of tags insertion. Note use UUID later */
	/* Insert the foreign keys to the tables */
	sqlStatement = `INSERT INTO news_tags(article_id, tag_id) VALUES`

	keyValues := []interface{}{}
	for i, tag_id := range tagIdList {
		keyValues = append(keyValues, articleId, tag_id)
		numFields := 2

		n := i * numFields
		sqlStatement += `(`
		for j := 0; j < numFields; j++ {
			sqlStatement += `$` + strconv.Itoa(n + j + 1) + `,`
		}
		sqlStatement = sqlStatement[:len(sqlStatement) - 1] + `),`
	}
	sqlStatement = sqlStatement[:len(sqlStatement) - 1]
	_, err = n.db.Exec(sqlStatement, keyValues...)
	if err != nil {
		return err
	}

	return nil
}

func (n newsArticleRepository) UpdateAnArticle(*NewsArticle) error {
	return nil
}
