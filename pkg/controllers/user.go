package controllers

import (
	"github.com/michaeltintiuc/shackle-api/pkg/services"
)

// newExpenseController creates a new ExpenseController
func newUserController(a controllable) {
	s := &services.User{Service: services.Service{Collection: a.Db().Collection("users")}, ClientCollection: a.Db().Collection("clientApplications")}

	a.AuthRouter().HandleFunc("/login", s.Login(a.Session())).Methods("POST")
	// a.ApiRouter().HandleFunc("/logout", s.Find(s.Model))).Methods("GET")

	a.ApiRouter().HandleFunc("/users", s.Find(&s.Model)).Methods("GET")
	a.ApiRouter().HandleFunc("/user/{id}", s.FindOne(&s.Model)).Methods("GET")
	a.ApiRouter().HandleFunc("/user", s.Create).Methods("POST")
	a.ApiRouter().HandleFunc("/user/{id}", s.Delete).Methods("DELETE")
}
