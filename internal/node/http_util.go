package node

import (
	"encoding/json"
	"net/http"
)

const CONTENT_TYPE = "application/json"

type errResp struct {
	Error string
}

func writeResponse(w http.ResponseWriter, content interface{}, httpStatusCode int) {
	jsonResp, err := json.Marshal(content)
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.WriteHeader(httpStatusCode)
	w.Write(jsonResp)
}

func writeErrorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	jsonError, _ := json.Marshal(errResp{err.Error()})
	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.WriteHeader(httpStatusCode)
	w.Write((jsonError))
}
