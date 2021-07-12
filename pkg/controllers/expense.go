package controllers

import (
	"github.com/michaeltintiuc/shackle-api/pkg/services"
)

// newExpenseController creates a new ExpenseController
func newExpenseController(a controllable) {
	s := &services.Expense{Service: services.Service{Collection: a.Db().Collection("expenses")}}
	a.ApiRouter().HandleFunc("/expenses", s.Find).Methods("GET")
	a.ApiRouter().HandleFunc("/expense/{id}", s.FindOne).Methods("GET")
	a.ApiRouter().HandleFunc("/expense", s.Create).Methods("POST")
	a.ApiRouter().HandleFunc("/expense/{id}", s.Delete).Methods("DELETE")
}
