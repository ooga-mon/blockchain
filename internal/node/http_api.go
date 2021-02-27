package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

const queryKeyIP = "ip"
const queryKeyPort = "port"
const queryKeyID = "id"

func (n *Node) handlerGetBlockChain(w http.ResponseWriter, r *http.Request) {
	err := writeResponse(w, n.db.Blocks, http.StatusOK)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("Retrieved blocks.")
}

func (n *Node) handlerMineBlock(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != CONTENT_TYPE {
		writeErrorResponse(w, fmt.Errorf("Content Type is not application/json"), http.StatusUnsupportedMediaType)
		return
	}

	var payload database.Transactions

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	block := n.Mine(payload)
	err = writeResponse(w, block, http.StatusOK)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("New block added.")

	n.postSync()
}

func (n *Node) handlerPostStatus(w http.ResponseWriter, r *http.Request) {
	peerStatus := Status{}
	err := readRequestBody(r, &peerStatus)
	if err != nil {
		return
	}

	n.updateStatusDifferences(peerStatus)

	nodeStatus := Status{n.info, n.peers, n.db}
	err = writeResponse(w, nodeStatus, http.StatusOK)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("sent diff in status.")
}
