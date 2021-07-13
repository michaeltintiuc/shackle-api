package controllers

import (
	"github.com/michaeltintiuc/shackle-api/pkg/models"
	"github.com/michaeltintiuc/shackle-api/pkg/services"
)

// newExpenseController creates a new ExpenseController
func newExpenseController(a controllable) {
	s := &services.Service{Collection: a.Db().Collection("expenses")}
	a.ApiRouter().HandleFunc("/expenses", s.Find(new(models.Expense))).Methods("GET")
	a.ApiRouter().HandleFunc("/expense/{id}", s.FindOne(new(models.Expense))).Methods("GET")
	a.ApiRouter().HandleFunc("/expense", s.Create(new(models.Expense))).Methods("POST")
	a.ApiRouter().HandleFunc("/expense/{id}", s.Delete).Methods("DELETE")
}
