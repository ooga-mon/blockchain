package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

type Server struct {
	Chain database.Blockchain
}

func (s *Server) GetBlockChain(w http.ResponseWriter, r *http.Request) {
	rsp, err := json.Marshal(s.Chain)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Retrieved blocks.")
	w.Write(rsp)
}
