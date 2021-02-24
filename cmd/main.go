package main

import (
	"log"
	"net/http"

	"github.com/ooga-mon/blockchain/app"
)

const HTTP_PORT = "8081"

var server app.Server

func handleRequests() {
	http.HandleFunc("/blocks", server.GetBlockChain)
	err := http.ListenAndServe(":"+HTTP_PORT, nil)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	server = app.NewServer()
	handleRequests()
}
