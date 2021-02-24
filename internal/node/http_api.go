package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

func (s *Node) getBlockChain(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, s.db.Blocks, http.StatusOK)
	fmt.Println("Retrieved blocks.")
}

func (s *Node) mineBlock(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeErrorResponse(w, errors.New("Content Type is not application/json"), http.StatusUnsupportedMediaType)
		return
	}

	var payload database.Payload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	block := s.db.AddBlock(payload)
	writeResponse(w, block, http.StatusOK)
	fmt.Println("New block added.")
}
