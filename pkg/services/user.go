package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/michaeltintiuc/shackle-api/pkg/models"
	"github.com/michaeltintiuc/shackle-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// User handles all item related logic
type User struct {
	Service
	Model models.User
}
type auth struct {
	Email    string
	Password string
}
type response struct {
	Token string `json:"token"`
}

// Login authenticates users
func (s *User) Login(jwtSecret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := auth{}
		json.NewDecoder(r.Body).Decode(&auth)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		cur := s.Collection.FindOne(ctx, bson.M{"email": auth.Email})

		if utils.HasError(w, cur.Err(), "Wrong email or password", http.StatusUnauthorized) {
			return
		}

		cur.Decode(&s.Model)

		err := bcrypt.CompareHashAndPassword([]byte(s.Model.Password), []byte(auth.Password))
		if utils.HasError(w, err, "Wrong email or password", http.StatusForbidden) {
			return
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &utils.JWTClaims{
			Uid:   s.Model.Id.String(),
			Email: s.Model.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),
			},
		}).SignedString([]byte(jwtSecret))

		if utils.HasError(w, err, "Failed token creation", http.StatusInternalServerError) {
			return
		}

		json.NewEncoder(w).Encode(response{token})
	}
}

func (s *User) Create(w http.ResponseWriter, r *http.Request) {
	// TODO
}
