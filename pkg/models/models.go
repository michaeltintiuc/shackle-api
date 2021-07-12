package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Expense model
type Expense struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Value float32            `json:"value,omitempty" bson:"value"`
}

// User model
type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"-" bson:"-"`
}
