package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func HasError(w http.ResponseWriter, err error, msg string, status int) bool {
	if err != nil {
		w.WriteHeader(status)
		log.Println(err)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{msg})
		return true
	}
	return false
}
