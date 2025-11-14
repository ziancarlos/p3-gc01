package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductID     primitive.ObjectID `json:"product_id" bson:"product_id" validate:"required"`
	Date          time.Time          `json:"date" bson:"date"`
	Price         float64            `json:"price" bson:"price" validate:"required,gt=0"`
	PaymentMethod string             `json:"payment_method" bson:"payment_method" validate:"required"`
	PaymentID     string             `json:"payment_id" bson:"payment_id"`
}

type TransactionRequest struct {
	ProductID     string  `json:"product_id" validate:"required"`
	Price         float64 `json:"price" validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	PaymentID     string  `json:"payment_id"`
}

type TransactionUpdateRequest struct {
	ProductID     string  `json:"product_id" validate:"omitempty"`
	Price         float64 `json:"price" validate:"omitempty,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"omitempty"`
	PaymentID     string  `json:"payment_id"`
}

type TransactionResponse struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	Date          time.Time `json:"date"`
	Price         float64   `json:"price"`
	PaymentMethod string    `json:"payment_method"`
	PaymentID     string    `json:"payment_id"`
}
