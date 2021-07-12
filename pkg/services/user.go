package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/michaeltintiuc/shackle-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User handles all item related logic
type User struct {
	Service
	Model models.User
}

// Create an item
func (s *User) Login(jwtSecret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := struct {
			Email    string
			Password string
		}{}
		json.NewDecoder(r.Body).Decode(&auth)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		cur := s.Collection.FindOne(ctx, bson.M{"email": auth.Email})
		model := models.User{}
		cur.Decode(&model)

		err := bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(auth.Password))
		if hasErrors(w, err, "Failed authentication", http.StatusForbidden) {
			return
		}

		token := jwt.New(jwt.SigningMethodHS256)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if hasErrors(w, err, "Failed token creation", http.StatusInternalServerError) {
			return
		}

		json.NewEncoder(w).Encode(struct{ Token string }{tokenString})
	}
}

// Create an item
func (s *User) Create(w http.ResponseWriter, r *http.Request) {
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
func (s *User) Find(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cur, err := s.Collection.Find(ctx, bson.M{})
	if hasErrors(w, err, "Failed fetching items", http.StatusInternalServerError) {
		return
	}

	var result []models.User
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
func (s *User) FindOne(w http.ResponseWriter, r *http.Request) {
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
func (s *User) Delete(w http.ResponseWriter, r *http.Request) {
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
