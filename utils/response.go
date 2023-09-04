package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adhupraba/breadit-server/internal/types"
)

type response struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error:", msg)
	}

	RespondWithJson(w, code, types.Json{"message": msg})
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	isError := false

	if code >= 400 && code < 600 {
		isError = true
	}

	jsonRes := response{Error: isError, Data: payload}

	RespondWithJsonDirect(w, code, jsonRes)
}

func RespondWithJsonDirect(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON response: %v\n", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
