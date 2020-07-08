package model

import (
	"io/ioutil"
	"log"
	"net/http"

	//"fmt"
	"encoding/json"

	"github.com/gorilla/mux"
)

type ContentManagerService struct {
	newsRepository NewsArticleRepository
}

func NewContentManagerService() (*ContentManagerService, error) {
	newsRepo, err := SetupDatabase()
	if err != nil {
		log.Fatalf("Error when attempting to setup the database")
	}

	return &ContentManagerService{newsRepository: newsRepo}, err
}

// TODO: RETURN THE UUID
func (c *ContentManagerService) CreateAnArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var articleData NewsArticleData
	json.Unmarshal(reqBody, &articleData)

	article := NewNewsArticle(articleData)
	err = c.newsRepository.CreateAnArticle(article)

	if err != nil {
		panic(err)
	}
}

func (c *ContentManagerService) UpdateAnArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	articleUUID := params["uuid"]

	/* Refactor this */
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var articleData NewsArticleData
	json.Unmarshal(reqBody, &articleData)

	article := NewNewsArticle(articleData)
	/* Till here */

	err = c.newsRepository.UpdateAnArticle(article, articleUUID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (c *ContentManagerService) RetrieveAnArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	articleUUID := params["uuid"]

	articleData, err := c.newsRepository.RetrieveAnArticle(articleUUID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusOK, articleData)
}

func (c *ContentManagerService) DeleteAnArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	articleUUID := params["uuid"]
	err := c.newsRepository.DeleteAnArticle(articleUUID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (c *ContentManagerService) SearchArticlesByTopic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	topic := params["topic"]
	article_list, err := c.newsRepository.FindNewsByTopic(topic)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusOK, article_list)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
