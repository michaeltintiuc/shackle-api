package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Text: %v\n", vars["text"])
}

func Init(r *mux.Router) {
	r.HandleFunc("/{text}", homePage)
	r.HandleFunc("/", homePage)
}
