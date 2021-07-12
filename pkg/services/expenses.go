package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/michaeltintiuc/shackle-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Expense handles all item related logic
type Expense struct {
	Service
	Model models.Expense
}

// Create an item
func (s *Expense) Create(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&s.Model)
	if hasErrors(w, err, "Failed to read item", http.StatusInternalServerError) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result, err := s.Collection.InsertOne(ctx, s.Model)

	if hasErrors(w, err, "Failed to create item", http.StatusInternalServerError) {
		return
	}

	json.NewEncoder(w).Encode(result)
}

// Find all items
func (s *Expense) Find(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cur, err := s.Collection.Find(ctx, bson.M{})
	if hasErrors(w, err, "Failed fetching items", http.StatusInternalServerError) {
		return
	}

	var result []models.Expense
	for cur.Next(context.Background()) {
		err := cur.Decode(&s.Model)
		if hasErrors(w, err, "Failed parsing items", http.StatusInternalServerError) {
			return
		}
		result = append(result, s.Model)
	}

	json.NewEncoder(w).Encode(result)
}

// FindOne item
func (s *Expense) FindOne(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := getId(r)
	if hasErrors(w, err, "Failed deleting item", http.StatusInternalServerError) {
		return
	}

	err = s.Collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&s.Model)

	if hasErrors(w, err, "Failed fetching item", http.StatusInternalServerError) {
		return
	}

	json.NewEncoder(w).Encode(s.Model)
}

// Delete an item by id
func (s *Expense) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := getId(r)
	if hasErrors(w, err, "Failed deleting item", http.StatusInternalServerError) {
		return
	}

	result, err := s.Collection.DeleteOne(ctx, bson.M{"_id": id})

	if hasErrors(w, err, "Failed deleting item", http.StatusInternalServerError) {
		return
	}

	json.NewEncoder(w).Encode(result)
}
