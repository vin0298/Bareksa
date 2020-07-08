package main

import (
	"fmt"
	//	"./model"
	"net/http"

	config "./config"
	"./routes"
	"github.com/spf13/viper"
)

func main() {
	config.SetupConfig(`.`)
	startServer()
}

func startServer() {
	router := routes.SetupRouter()
	err := http.ListenAndServe(viper.GetString("SERVER_PORT"), router)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("server is running")
}
