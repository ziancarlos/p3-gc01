package repository

import (
	"context"
	"p3-graded-challenge-1-ziancarlos/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
}

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(db *mongo.Database) PaymentRepository {
	return &paymentRepository{
		collection: db.Collection("payments"),
	}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return err
	}
	payment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
