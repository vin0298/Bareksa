package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type NewsArticleRepository interface {
	FindNewsByTopic(topic string) ([]ArticleReadModel, error)
	FindNewsByStatus(status string) ([]ArticleReadModel, error)
	RetrieveAnArticle(uuid string) (*ArticleReadModel, error)
	DeleteAnArticle(articleID string) error
	CreateAnArticle(article *NewsArticle) error
	UpdateAnArticle(article *NewsArticle, articleUUID string) error
}

type newsArticleRepository struct {
	db *sql.DB
}

/* TEMP */
func SetupDatabase() (NewsArticleRepository, error) {
	/* Enable SSL later */
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"), viper.GetInt("DB_PORT"),
		viper.GetString("DB_USER"), viper.GetString("DB_PW"),
		viper.GetString("DB_NAME"))

	postgresDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("%s", err)
		return newsArticleRepository{}, err
	}

	err = postgresDB.Ping()
	if err != nil {
		log.Printf("%s", err)
		return newsArticleRepository{}, err
	}

	return newsArticleRepository{db: postgresDB}, nil
}

func (n newsArticleRepository) FindNewsByTopic(topic string) ([]ArticleReadModel, error) {
	/* Retrieve the news_articles */
	/**
		SELECT news_articles.*, string_agg(tag_name, ',') as tag_names FROM news_articles
	JOIN news_tags ON news_tags.article_id=news_articles.article_id
	JOIN tags ON news_tags.tag_id=tags.tag_id GROUP BY news_articles.article_id;
	*/
	sqlStatement := `SELECT news_articles.*, string_agg(tag_name, ',') as tag_names
						FROM news_articles 
						JOIN news_tags ON news_tags.article_id=news_articles.article_id
						JOIN tags ON news_tags.tag_id=tags.tag_id
						WHERE news_articles.topic=$1
						GROUP BY news_articles.article_id;`

	rows, err := n.db.Query(sqlStatement, topic)
	var listOfArticles []ArticleReadModel
	if err != nil {
		log.Printf("Error at FindNewByTopic() on query: %s", err)
		return listOfArticles, err
	}

	defer rows.Close()

	for rows.Next() {
		article := ArticleReadModel{}
		articlePK := 0
		var tag_names string
		err = rows.Scan(&articlePK, &article.Author, &article.Title,
			&article.Content, &article.Time_published, &article.Uuid,
			&article.Topic, &article.Status, &tag_names)

		article.Tags = strings.Split(tag_names, `,`)
		listOfArticles = append(listOfArticles, article)
	}

	err = rows.Err()
	if err != nil {
		return listOfArticles, err
	}

	return listOfArticles, nil
}

/* Can be refactor for DRY */
func (n newsArticleRepository) FindNewsByStatus(status string) ([]ArticleReadModel, error) {
	sqlStatement := `SELECT news_articles.*, string_agg(tag_name, ',') as tag_names
						FROM news_articles 
						JOIN news_tags ON news_tags.article_id=news_articles.article_id
						JOIN tags ON news_tags.tag_id=tags.tag_id
						WHERE news_articles.status=$1
						GROUP BY news_articles.article_id;`

	rows, err := n.db.Query(sqlStatement, status)
	var listOfArticles []ArticleReadModel
	if err != nil {
		log.Printf("Error at FindNewByTopic() on query: %s", err)
		return listOfArticles, err
	}

	defer rows.Close()

	for rows.Next() {
		article := ArticleReadModel{}
		articlePK := 0
		var tag_names string
		err = rows.Scan(&articlePK, &article.Author, &article.Title,
			&article.Content, &article.Time_published, &article.Uuid,
			&article.Topic, &article.Status, &tag_names)

		article.Tags = strings.Split(tag_names, `,`)
		listOfArticles = append(listOfArticles, article)
	}

	err = rows.Err()
	if err != nil {
		return listOfArticles, err
	}

	return listOfArticles, nil
}

func (n newsArticleRepository) RetrieveAnArticle(uuid string) (*ArticleReadModel, error) {
	sqlStatement := `SELECT * FROM news_articles WHERE uuid=$1;`
	row := n.db.QueryRow(sqlStatement, uuid)

	var newsModel = ArticleReadModel{}
	var articlePK int

	err := row.Scan(&articlePK, &newsModel.Author,
		&newsModel.Title, &newsModel.Content,
		&newsModel.Time_published, &newsModel.Uuid,
		&newsModel.Topic, &newsModel.Status)

	if err != nil {
		return &ArticleReadModel{}, err
	}

	/* Retrieve all the tags related to it */
	sqlStatement = `SELECT string_agg(tag_name, ',') FROM news_tags JOIN tags 
						ON news_tags.tag_id=tags.tag_id 
						WHERE news_tags.article_id=$1`

	var tagListStr string
	err = n.db.QueryRow(sqlStatement, articlePK).Scan(&tagListStr)

	if err != nil {
		return &ArticleReadModel{}, err
	}
	newsModel.Tags = strings.Split(tagListStr, `,`)
	/* Done tag retrieval */

	return &newsModel, nil
}

func (n newsArticleRepository) DeleteAnArticle(articleID string) error {
	sqlStatement := `DELETE FROM news_articles WHERE uuid=$1;`
	_, err := n.db.Exec(sqlStatement, articleID)
	if err != nil {
		log.Printf("Error when deleting a NewsArticle: %s", err)
		return err
	}

	return nil
}

func (n newsArticleRepository) CreateAnArticle(article *NewsArticle) error {
	sqlStatement := `INSERT INTO news_articles(author, title, content, time_published, uuid, topic, status)
					 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING article_id`

	/* Insert the NewsArticle content */
	articleID := 0
	err := n.db.QueryRow(sqlStatement, article.Author(), article.Title(),
		article.Content(), article.TimePublished(), article.Id(),
		article.GetTopic(), article.Status()).Scan(&articleID)

	if err != nil {
		log.Printf("Error when inserting a new NewsArticle: %s", err)
		return err
	}

	/* Refactor: Insert the tags */
	tagIdList, err := n.insertArticleTags(article)
	if err != nil {
		return err
	}

	/* End of tags insertion. Note use UUID later */
	/* Insert the foreign keys to the tables */
	err = n.updateArticleTagsJoinTable(articleID, tagIdList)
	if err != nil {
		return err
	}

	return nil
}

func (n newsArticleRepository) UpdateAnArticle(article *NewsArticle, articleUUID string) error {
	sqlStatement := `UPDATE news_articles SET author=$1, title=$2,
					content=$3, time_published=$4, topic=$5, status=$6
					WHERE uuid=$7 RETURNING article_id;`

	articleID := 0
	err := n.db.QueryRow(sqlStatement, article.Author(), article.Title(), article.Content(),
		article.TimePublished(), article.GetTopic(), article.Status(), articleUUID).Scan(&articleID)

	if articleID == 0 || err != nil {
		if err == nil {
			return errors.New("Invalid update data")
		}

		return err
	}

	tagIdList, err := n.insertArticleTags(article)
	if err != nil {
		return err
	}

	err = n.updateArticleTagsJoinTable(articleID, tagIdList)
	if err != nil {
		return err
	}

	var tagNameList []string
	for _, tagObj := range article.Tags() {
		tagNameList = append(tagNameList, tagObj.Name())
	}

	err = n.deleteUnusedTagsJoinTable(articleID, tagNameList)
	if err != nil {
		return err
	}

	return nil
}

func (n newsArticleRepository) insertArticleTags(article *NewsArticle) ([]int, error) {
	sqlStatement := `INSERT INTO tags(tag_name, uuid) VALUES`
	tagValues := []interface{}{}
	for i, tag := range article.Tags() {
		tagValues = append(tagValues, tag.Name(), tag.Id())
		numFields := 2

		n := i * numFields
		sqlStatement += `(`
		for j := 0; j < numFields; j++ {
			sqlStatement += `$` + strconv.Itoa(n+j+1) + `,`
		}
		sqlStatement = sqlStatement[:len(sqlStatement)-1] + `),`
	}
	sqlStatement = sqlStatement[:len(sqlStatement)-1]
	sqlStatement += `ON CONFLICT DO NOTHING RETURNING tag_id`

	rows, err := n.db.Query(sqlStatement, tagValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tagIdList []int

	for rows.Next() {
		var tag_id int
		err = rows.Scan(&tag_id)
		if err != nil {
			return tagIdList, err
		}
		tagIdList = append(tagIdList, tag_id)
	}

	err = rows.Err()
	if err != nil {
		return tagIdList, err
	}

	return tagIdList, nil
}

func (n newsArticleRepository) updateArticleTagsJoinTable(articleID int, tagIdList []int) error {
	if len(tagIdList) == 0 {
		return nil
	}

	sqlStatement := `INSERT INTO news_tags(article_id, tag_id) VALUES`

	keyValues := []interface{}{}
	for i, tag_id := range tagIdList {
		keyValues = append(keyValues, articleID, tag_id)
		numFields := 2

		n := i * numFields
		sqlStatement += `(`
		for j := 0; j < numFields; j++ {
			sqlStatement += `$` + strconv.Itoa(n+j+1) + `,`
		}
		sqlStatement = sqlStatement[:len(sqlStatement)-1] + `),`
	}
	sqlStatement = sqlStatement[:len(sqlStatement)-1]
	log.Printf("SQL STATEMENT: %s", sqlStatement)
	_, err := n.db.Query(sqlStatement, keyValues...)
	if err != nil {
		log.Printf("Encountered an error when updating the join table")
		return err
	}

	return nil
}

func (n newsArticleRepository) findAllTagsOfAnArticle(articleUUID uuid.UUID) ([]TagReadModel, error) {
	sqlStatement := `SELECT tags.uuid, tags.tag_name, tags.tag_id FROM news_articles 
					JOIN news_tags ON news_articles.article_id=news_tags.article_id 
					JOIN tags ON tags.tag_id=news_tags.tag_id
					WHERE news_articles.uuid=$1;`

	var tagList []TagReadModel
	rows, err := n.db.Query(sqlStatement, articleUUID)
	if err != nil {
		log.Printf("Error at findAllTagsOfAnArticle(): %s", err)
		return tagList, err
	}

	defer rows.Close()
	for rows.Next() {
		var tagObj TagReadModel
		err = rows.Scan(&tagObj.Uuid, &tagObj.Name, &tagObj.Id)
		if err != nil {
			return tagList, err
		}
		tagList = append(tagList, tagObj)
	}

	err = rows.Err()
	if err != nil {
		return tagList, err
	}

	return tagList, nil
}

func (n newsArticleRepository) deleteUnusedTagsJoinTable(articleID int, tagList []string) error {
	sqlStatement := `SELECT news_tags.tag_id FROM news_tags 
						JOIN tags ON news_tags.tag_id=tags.tag_id
						WHERE NOT (tags.tag_name = ANY(array[`

	tagNameList := []interface{}{}
	for i, tagName := range tagList {
		log.Printf("tagname: %s", tagName)
		tagNameList = append(tagNameList, tagName)
		sqlStatement += `$` + strconv.Itoa(i+1) + `,`
	}
	sqlStatement = sqlStatement[:len(sqlStatement)-1]
	sqlStatement += fmt.Sprintf(`])) AND news_tags.article_id=%d`, articleID)

	log.Printf("delete sql: %s", sqlStatement)
	rows, err := n.db.Query(sqlStatement, tagNameList...)
	if err != nil {
		log.Printf("Initial filter failed")
		return err
	}

	defer rows.Close()
	var tagIdToRemove []int

	for rows.Next() {
		var tag_id int
		err = rows.Scan(&tag_id)
		if err != nil {
			log.Printf("Error when scanning the tag id at DeleteUnusedTagsJoinTable(): %s", err)
			return err
		}
		tagIdToRemove = append(tagIdToRemove, tag_id)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	sqlStatement = `DELETE FROM news_tags WHERE news_tags.article_id=$1 
						AND news_tags.tag_id IN (`

	tagIdList := []interface{}{}
	tagIdList = append(tagIdList, articleID)
	for i, id := range tagIdToRemove {
		tagIdList = append(tagIdList, id)
		sqlStatement += `$` + strconv.Itoa(i+2) + `,`
	}

	sqlStatement = sqlStatement[:len(sqlStatement)-1]
	sqlStatement += `)`

	_, err = n.db.Exec(sqlStatement, tagIdList...)
	if err != nil {
		log.Printf("Error when deleting unused tags: %s", err)
		return err
	}

	return nil
}
