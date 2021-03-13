package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/octothorped/blockchain/database"
)

const queryKeyIP = "ip"
const queryKeyPort = "port"
const queryKeyID = "id"

func (n *Node) handlerGetBlockChain(w http.ResponseWriter, r *http.Request) {
	err := WriteResponse(w, n.db.Blocks, http.StatusOK)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("Retrieved blocks.")
}

func (n *Node) handlerPostTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != CONTENT_TYPE {
		WriteErrorResponse(w, fmt.Errorf("Content Type is not application/json"), http.StatusUnsupportedMediaType)
		return
	}

	var payload database.SignedTransaction

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = n.addPendingTransaction(payload)
	if err != nil {
		WriteErrorResponse(w, fmt.Errorf("Could not authenticate transaction"), http.StatusBadRequest)
		return
	}

	n.broadcastTransaction(payload)
}

func (n *Node) handlerPostSyncTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != CONTENT_TYPE {
		WriteErrorResponse(w, fmt.Errorf("Content Type is not application/json"), http.StatusUnsupportedMediaType)
		return
	}

	var payload database.SignedTransaction

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = n.addPendingTransaction(payload)
}

func (n *Node) handlerPostStatus(w http.ResponseWriter, r *http.Request) {
	peerStatus := Status{}
	err := ReadRequestBody(r, &peerStatus)
	if err != nil {
		return
	}

	n.updateStatusDifferences(peerStatus)

	nodeStatus := Status{n.info, n.peers, n.db}
	err = WriteResponse(w, nodeStatus, http.StatusOK)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("sent diff in status.")
}
