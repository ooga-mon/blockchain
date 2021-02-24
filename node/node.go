package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

type Server struct {
	Chain database.Blockchain
}

func NewServer() Server {
	return Server{database.NewBlockchain()}
}

func (s *Server) GetBlockChain(w http.ResponseWriter, r *http.Request) {
	rsp, err := json.Marshal(s.Chain)
	if err != nil {
		fmt.Print(err)
		writeResponseMessage(w, "Error retrieving blockchain.", http.StatusInternalServerError)
		return
	}
	fmt.Println("Retrieved blocks.")
	w.Write(rsp)
}

func (s *Server) MineBlock(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeResponseMessage(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var payload database.Payload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		writeResponseMessage(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	block := s.Chain.AddBlock(payload)

	blockJson, err := json.Marshal(block)
	if err != nil {
		fmt.Print(err)
		writeResponseMessage(w, "Block added to chain. Error creating new block response. Check using GET block request.", http.StatusInternalServerError)
		return
	}

	fmt.Println("New block added.")
	w.Write(blockJson)
}

func writeResponseMessage(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
