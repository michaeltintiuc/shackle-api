package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseExpense struct {
	Id                primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string               `json:"name" bson:"name"`
	Value             float32              `json:"value" bson:"value"`
	Currency          string               `json:"currency" bson:"currency"`
	Discount          float32              `json:"discount" bson:"discount"`
	DiscountIsPercent bool                 `json:"discountIsPercent" bson:"discountIsPercent"`
	Tags              []primitive.ObjectID `json:"tags" bson:"tags"`
	Receipt           string               `json:"receipt" bson:"receipt"`
	CreatedAt         time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time            `json:"updatedAt" bson:"updatedAt"`
}

// Expense model
type Expense struct {
	Id                primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string               `json:"name" bson:"name"`
	Value             float32              `json:"value" bson:"value"`
	Currency          string               `json:"currency" bson:"currency"`
	Discount          float32              `json:"discount" bson:"discount"`
	DiscountIsPercent bool                 `json:"discountIsPercent" bson:"discountIsPercent"`
	Tags              []primitive.ObjectID `json:"tags" bson:"tags"`
	Receipt           string               `json:"receipt" bson:"receipt"`
	CreatedAt         time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time            `json:"updatedAt" bson:"updatedAt"`
	Recurring         bool                 `json:"recurring" bson:"recurring"`
	Items             []BaseExpense        `json:"items" bson:"items"`
}

// User model
type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"-" bson:"password"`
}

type ClientApp struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Type string             `json:"type" bson:"type"`
	Name string             `json:"name" bson:"name"`
}
