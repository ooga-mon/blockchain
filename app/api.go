package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const HTTP_PORT = "8081"

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var articles = []Article{
	{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func ReturnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(articles)
}
