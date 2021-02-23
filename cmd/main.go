package main

import (
	"log"
	"net/http"

	"blockchain/app/api"
)

const HTTP_PORT = "8081"

func handleRequests() {
	http.HandleFunc("/", api.HomePage)
	http.HandleFunc("/articles", api.ReturnAllArticles)
	err := http.ListenAndServe(":"+HTTP_PORT, nil)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	handleRequests()
}
