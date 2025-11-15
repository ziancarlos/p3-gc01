package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"p3-graded-challenge-1-ziancarlos/config"
	"p3-graded-challenge-1-ziancarlos/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindAll() ([]models.Transaction, error)
	FindByID(id primitive.ObjectID) (*models.Transaction, error)
	Update(id primitive.ObjectID, update bson.M) error
	Delete(id primitive.ObjectID) error
}

type transactionRepository struct {
	collection *mongo.Collection
	cfg        *config.Config
}

type PaymentRequest struct {
	Amount float64 `json:"amount"`
}

type PaymentResponse struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

func NewTransactionRepository(db *mongo.Database, cfg *config.Config) TransactionRepository {
	return &transactionRepository{
		collection: db.Collection("transactions"),
		cfg:        cfg,
	}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// HTTP call ke payment service
	paymentReq := PaymentRequest{
		Amount: transaction.Price,
	}

	jsonData, err := json.Marshal(paymentReq)
	if err != nil {
		return fmt.Errorf("failed to marshal payment request: %w", err)
	}

	paymentURL := fmt.Sprintf("%s/payments", r.cfg.PaymentService.BaseURI)
	req, err := http.NewRequestWithContext(ctx, "POST", paymentURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create payment request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call payment service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("payment service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read payment response: %w", err)
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(bodyBytes, &paymentResp); err != nil {
		return fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	// Simpan payment ID ke transaction
	transaction.PaymentID = paymentResp.ID
	transaction.Date = time.Now()

	result, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		return err
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *transactionRepository) FindAll() ([]models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) FindByID(id primitive.ObjectID) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction models.Transaction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) Update(id primitive.ObjectID, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateDoc := bson.M{"$set": update}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *transactionRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

