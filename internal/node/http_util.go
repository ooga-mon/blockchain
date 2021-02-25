package node

import (
	"encoding/json"
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
