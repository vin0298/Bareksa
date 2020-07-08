package routes

import (
	"log"

	"../services"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	contentManager, err := services.NewContentManagerService()
	if err != nil {
		log.Fatalf("Creating the ContentManagerService Failed %s", err)
	}

	// TODO: ADD ERROR FOR UNHANDLED ROUTES

	router.HandleFunc("/articles", contentManager.CreateAnArticle).Methods("POST")
	router.HandleFunc("/articles/{uuid}", contentManager.RetrieveAnArticle).Methods("GET")
	router.HandleFunc("/articles/{uuid}", contentManager.DeleteAnArticle).Methods("DELETE")
	router.HandleFunc("/articles/{uuid}", contentManager.UpdateAnArticle).Methods("PUT")

	router.HandleFunc("/articles/search-by-topic/{topic}", contentManager.SearchArticlesByTopic).Methods("GET")
	router.HandleFunc("/articles/search-by-status/{status}", contentManager.SearchArticlesByStatus).Methods("GET")
	router.HandleFunc("/articles", contentManager.ListAllArticles).Methods("GET")

	router.HandleFunc("/articles/tags", contentManager.CreateATag).Methods("POST")
	router.HandleFunc("/articles/tags/{uuid}", contentManager.RenameATag).Methods("PUT")
	router.HandleFunc("/articles/tags/{uuid}", contentManager.RetrieveATag).Methods("GET")
	router.HandleFunc("/articles/all/tags", contentManager.ListAllTags).Methods("GET")

	/* attach a tag to an article, remove a tag from an article
	router.HandleFunc("/articles/{article-uuid}/tags/{tag-uuid}",)
	*/

	return router
}
