package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Expense model
type Expense struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Value float32            `json:"value,omitempty" bson:"value,omitempty"`
}

// User model
type User struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"-" bson:"password"`
	SetPassword string             `json:"password" bson:"-"`
}
