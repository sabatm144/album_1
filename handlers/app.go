package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

type result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


func (e result) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e result) Success() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func renderJSON(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("ERROR: renderJson - %q\n", err)
	}
}

func renderERROR(w http.ResponseWriter, res *result) {
	renderJSON(w, res.Code, res)
}


