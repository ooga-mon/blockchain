package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const CONTENT_TYPE = "application/json"

type StandardResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"Error"`
}

type errResp struct {
	Error string `json:"Error"`
}

type Status struct {
	KnowPeers map[string]connectionInfo `json:"connection_info"`
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

func readResponse(r *http.Response, reqBody interface{}) error {
	resBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body. %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to process response. %s", string(resBodyJson))
	}

	err = json.Unmarshal(resBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response body. %s", err.Error())
	}

	return nil
}
