package main

import (
	"log"

	"net/http"

	config "github.com/bareksa/config"
	"github.com/bareksa/routes"
	"github.com/spf13/viper"
)

func main() {
	config.SetupConfig(`.`)
	startServer()
}

func startServer() {
	router := routes.SetupRouter()
	log.Printf("Server is running")
	log.Fatal(http.ListenAndServe(viper.GetString("SERVER_PORT"), router))
}
