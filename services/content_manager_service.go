package services

import (
	"io/ioutil"
	"log"
	"net/http"

	//"fmt"
	"encoding/json"

	model "github.com/bareksa/model"
	repo "github.com/bareksa/repository"

	"github.com/gorilla/mux"
)

// This is the domain service. The concept is working as the
// software of a Content-Management-Software such as Drupal/WordPress.
// The user will then "order" the CMS through the API. This "CMS"
// will then try to read the data provided, create the suitable
// aggregate and pass it to the repository for persistence

type ContentManagerService struct {
	newsRepository repo.NewsArticleRepository
}

func NewContentManagerService() (*ContentManagerService, error) {
	newsRepo, err := repo.SetupDatabase()
	if err != nil {
		log.Fatalf("Error when attempting to setup the database")
	}

	return &ContentManagerService{newsRepository: newsRepo}, err
}

// TODO: RETURN THE UUID
func (c *ContentManagerService) CreateAnArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var articleData model.NewsArticleData
	json.Unmarshal(reqBody, &articleData)

	article := model.NewNewsArticle(articleData)
	articleUUID, err := c.newsRepository.CreateAnArticle(article)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	res := struct{ uuid string }{articleUUID}

	RespondWithJSON(w, http.StatusCreated, res)
}

func (c *ContentManagerService) UpdateAnArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	articleUUID := params["uuid"]

	/* Refactor this */
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var articleData model.NewsArticleData
	json.Unmarshal(reqBody, &articleData)

	article := model.NewNewsArticle(articleData)
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

func (c *ContentManagerService) SearchArticlesByStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	status := params["status"]
	article_list, err := c.newsRepository.FindNewsByStatus(status)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusOK, article_list)
}

func (c *ContentManagerService) ListAllArticles(w http.ResponseWriter, r *http.Request) {
	article_list, err := c.newsRepository.RetrieveAllArticles()

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusOK, article_list)
}

/** Start of Tags Management Service **/
func (c *ContentManagerService) CreateATag(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Something went wrong when reading the body: %s", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var tagData model.TagData
	json.Unmarshal(reqBody, &tagData)

	newTag := model.CreateArticleTag(tagData.Name)
	tagResult, err := c.newsRepository.CreateATag(&newTag)

	if err != nil {
		log.Printf("Something went wrong when creating the tag: %s", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusCreated, tagResult)
}

func (c *ContentManagerService) RenameATag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tagUUID := params["uuid"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Something went wrong when reading the body: %s", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var tagData model.TagData
	json.Unmarshal(reqBody, &tagData)
	err = c.newsRepository.RenameATag(tagUUID, tagData.Name)
	if err != nil {
		log.Printf("Error occured when renaming the tag: %s", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (c *ContentManagerService) RetrieveATag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tagUUID := params["uuid"]

	tagData, err := c.newsRepository.SearchATag(tagUUID)
	if err != nil {
		log.Printf("Error when searching for a specific tag: %s", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondWithJSON(w, http.StatusOK, tagData)
}

func (c *ContentManagerService) ListAllTags(w http.ResponseWriter, r *http.Request) {
	tags_list, err := c.newsRepository.RetrieveAllTags()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	RespondWithJSON(w, http.StatusOK, tags_list)
}

// The two functions below are helper functions help to write the API
// response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
