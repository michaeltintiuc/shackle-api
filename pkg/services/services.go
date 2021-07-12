package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Collection *mongo.Collection
	Model      interface{}
}

func hasErrors(w http.ResponseWriter, err error, msg string, status int) bool {
	if err != nil {
		w.WriteHeader(status)
		log.Println(err)
		json.NewEncoder(w).Encode(struct{ Message string }{msg})
		return true
	}
	return false
}

func getId(r *http.Request) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(mux.Vars(r)["id"])
}
