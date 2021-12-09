package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/michaeltintiuc/shackle-api/pkg/models"
	"github.com/michaeltintiuc/shackle-api/pkg/session"
	"github.com/michaeltintiuc/shackle-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User handles all item related logic
type User struct {
	Service
	ClientCollection *mongo.Collection
	Model            models.User
}
type auth struct {
	Email    string
	Password string
	ClientId primitive.ObjectID
}
type response struct {
	Success bool `json:"success"`
}

// Login authenticates users
func (s *User) Login(sessionInfo *session.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := auth{}

		json.NewDecoder(r.Body).Decode(&auth)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// Find the user
		err := s.Collection.FindOne(ctx, bson.M{"email": auth.Email}).Decode(&s.Model)
		if utils.HasError(w, err, "Wrong email or password", http.StatusUnauthorized) {
			return
		}

		// Validate credentials
		err = bcrypt.CompareHashAndPassword([]byte(s.Model.Password), []byte(auth.Password))
		if utils.HasError(w, err, "Wrong email or password", http.StatusForbidden) {
			return
		}

		session, err := sessionInfo.Get(r)
		if utils.HasError(w, err, "Application error", http.StatusInternalServerError) {
			return
		}

		session.Values["id"] = s.Model.Id.Hex()

		err = session.Save(r, w)
		if utils.HasError(w, err, "Application error", http.StatusInternalServerError) {
			return
		}

		json.NewEncoder(w).Encode(response{Success: true})
	}
}

func (s *User) Create(w http.ResponseWriter, r *http.Request) {
	// TODO
}
