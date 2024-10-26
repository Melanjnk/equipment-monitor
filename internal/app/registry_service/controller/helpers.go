package controller

import (
	"fmt"
	"encoding/json"
	"net/http"
)

const parameterIsRequired string = "Parameter `%s` is required."

func writeMessage(w http.ResponseWriter, status int, message string, params ...interface{}) {
	w.WriteHeader(status)
	fmt.Fprintf(w, message, params...)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func readJSON(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}
