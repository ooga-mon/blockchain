package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

const CONTENT_TYPE = "application/json"

type StandardResponse struct {
	Success bool   `json:"Success"`
	Error   string `json:"Error"`
}

type errResp struct {
	Error string `json:"Error"`
}

type Status struct {
	Info       connectionInfo            `json:"Info"`
	KnowPeers  map[string]connectionInfo `json:"Known_Peers"`
	Blockchain database.Blockchain       `json:"Blockchain"` //TODO switch this to be the last block info instead of the whole chain each time
}

func writeRequest(url string, content interface{}) (*http.Response, error) {
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewBuffer(jsonContent))
}

func readRequestBody(r *http.Request, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body. %s", err.Error())
	}

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil
}

func writeResponse(w http.ResponseWriter, content interface{}, httpStatusCode int) error {
	jsonResp, err := json.Marshal(content)
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.WriteHeader(httpStatusCode)
	w.Write(jsonResp)
	return nil
}

func writeErrorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	jsonError, _ := json.Marshal(errResp{err.Error()})
	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.WriteHeader(httpStatusCode)
	w.Write((jsonError))
}

func readResponse(r *http.Response, respBody interface{}) error {
	resBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body. %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to process response. %s", string(resBodyJson))
	}

	err = json.Unmarshal(resBodyJson, respBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response body. %s", err.Error())
	}

	return nil
}
