package controllers

import (
	"github.com/michaeltintiuc/shackle-api/pkg/services"
)

// newExpenseController creates a new ExpenseController
func newUserController(a controllable, jwtSecret string) {
	u := &services.User{Service: services.Service{Collection: a.Db().Collection("users")}}

	a.AuthRouter().HandleFunc("/login", u.Login(jwtSecret)).Methods("POST")
	a.ApiRouter().HandleFunc("/logout", u.Find).Methods("GET")

	a.ApiRouter().HandleFunc("/users", u.Find).Methods("GET")
	a.ApiRouter().HandleFunc("/user", u.Create).Methods("POST")
	a.ApiRouter().HandleFunc("/user/{id}", u.Delete).Methods("DELETE")
}
