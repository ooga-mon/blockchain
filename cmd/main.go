package main

import (
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/node"
)

type flagParams struct {
}

const HTTP_PORT = "8081"

var server node.Server

func serverHttp() {
	http.HandleFunc("/blocks", server.GetBlockChain)
	http.HandleFunc("/mine", server.MineBlock)
	err := http.ListenAndServe(":"+HTTP_PORT, nil)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Listening for requests on port: %s.", HTTP_PORT)
}

func main() {
	server = node.NewServer()
	serverHttp()
}
