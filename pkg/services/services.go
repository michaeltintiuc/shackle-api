package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/shackle-api/pkg/models"
	"github.com/michaeltintiuc/shackle-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	Collection *mongo.Collection
}

// Create an item
func (s *Service) Create(model interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(model)
		if utils.HasError(w, err, "Failed to read item", http.StatusInternalServerError) {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		result, err := s.Collection.InsertOne(ctx, model)
		if utils.HasError(w, err, "Failed to create item", http.StatusInternalServerError) {
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

// Find all items
func (s *Service) Find(model interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		opts := options.Find()
		filter := bson.M{}
		parsedParams, err := parseQuery(r, map[string]string{"limit": "int", "search": "string", "tags": "uintArray"})

		if utils.HasError(w, err, "Failed parsing query parameters", http.StatusBadRequest) {
			return
		}

		for key, val := range parsedParams {
			switch key {
			case "limit":
				opts.SetLimit(val.(int64))
			case "search":
				filter["$text"] = bson.M{"$search": val}
				opts.SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})
			}
		}

		cur, err := s.Collection.Find(ctx, filter, opts)
		defer cur.Close(ctx)

		if utils.HasError(w, err, "Failed fetching items", http.StatusInternalServerError) {
			return
		}

		var result []interface{}
		for cur.Next(context.Background()) {
			err := cur.Decode(model)
			if utils.HasError(w, err, "Failed parsing items", http.StatusInternalServerError) {
				return
			}
			switch t := model.(type) {
			default:
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("unexpected type %T\n", t)
				json.NewEncoder(w).Encode(struct{ Message string }{"Failed parsing items"})
				return
			case *models.User:
				result = append(result, *model.(*models.User))
			case *models.Expense:
				result = append(result, *model.(*models.Expense))
			}
		}

		json.NewEncoder(w).Encode(result)
	}
}

// FindOne item
func (s *Service) FindOne(model interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
		if utils.HasError(w, err, "Failed deleting item", http.StatusInternalServerError) {
			return
		}

		err = s.Collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(model)
		if utils.HasError(w, err, "Failed fetching item", http.StatusInternalServerError) {
			return
		}

		json.NewEncoder(w).Encode(model)
	}
}

// Delete an item by id
func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if utils.HasError(w, err, "Failed deleting item", http.StatusInternalServerError) {
		return
	}

	result, err := s.Collection.DeleteOne(ctx, bson.M{"_id": id})

	if utils.HasError(w, err, "Failed deleting item", http.StatusInternalServerError) {
		return
	}

	json.NewEncoder(w).Encode(result)
}
