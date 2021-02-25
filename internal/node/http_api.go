package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

	block := n.db.AddBlock(payload)
	err = writeResponse(w, block, http.StatusOK)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("New block added.")
}

func (n *Node) handlerAddPeer(w http.ResponseWriter, r *http.Request) {
	peerIP := r.URL.Query().Get(queryKeyIP)
	peerPortStr := r.URL.Query().Get(queryKeyPort)

	peerPort, err := strconv.ParseUint(peerPortStr, 10, 32)
	if err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	peerCon := newConnectionInfo(peerIP, peerPort, true)

	n.addPeer(peerCon)

	err = writeResponse(w, StandardResponse{true, ""}, http.StatusOK)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("New peer added.")
}
