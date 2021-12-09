package controllers

import (
	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/shackle-api/pkg/session"

	"go.mongodb.org/mongo-driver/mongo"
)

type controllable interface {
	Db() *mongo.Database
	ApiRouter() *mux.Router
	AuthRouter() *mux.Router
	Session() *session.Session
}

// Init initializes all of the controllers
func Init(a controllable) {
	newExpenseController(a)
	newUserController(a)
	// a.Router().HandleFunc("/ping", router.PingDb)
	// a.Router().HandleFunc("/databases", router.ListDatabases)
	// a.Router().HandleFunc("/{text}", router.homePage)
	// a.Router().HandleFunc("/", router.homePage)
}
