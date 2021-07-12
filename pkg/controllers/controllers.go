package controllers

import (
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
)

type controllable interface {
	Db() *mongo.Database
	ApiRouter() *mux.Router
	AuthRouter() *mux.Router
}

// Init initializes all of the controllers
func Init(a controllable, jwtSecret string) {
	newExpenseController(a)
	newUserController(a, jwtSecret)
	// a.Router().HandleFunc("/ping", router.PingDb)
	// a.Router().HandleFunc("/databases", router.ListDatabases)
	// a.Router().HandleFunc("/{text}", router.homePage)
	// a.Router().HandleFunc("/", router.homePage)
}
