package service

import (
	"errors"
	"shopping-service/config"
	"shopping-service/models"
	"shopping-service/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionService interface {
	CreateTransaction(req *models.TransactionRequest) (*models.TransactionResponse, error)
	GetAllTransactions() ([]models.TransactionResponse, error)
	GetTransactionByID(id string) (*models.TransactionResponse, error)
	UpdateTransaction(id string, req *models.TransactionUpdateRequest) error
	DeleteTransaction(id string) error
}

type transactionService struct {
	repo repository.TransactionRepository
	cfg  *config.Config
}

func NewTransactionService(repo repository.TransactionRepository, cfg *config.Config) TransactionService {
	return &transactionService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *transactionService) CreateTransaction(req *models.TransactionRequest) (*models.TransactionResponse, error) {
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return nil, errors.New("invalid product_id")
	}
	transaction := &models.Transaction{
		ProductID:     productID,
		Price:         req.Price,
		PaymentMethod: req.PaymentMethod,
	}

	if err := s.repo.Create(transaction); err != nil {
		return nil, err
	}

	response := &models.TransactionResponse{
		ID:            transaction.ID.Hex(),
		ProductID:     transaction.ProductID.Hex(),
		Date:          transaction.Date,
		Price:         transaction.Price,
		PaymentMethod: transaction.PaymentMethod,
		PaymentID:     transaction.PaymentID,
	}

	return response, nil
}

func (s *transactionService) GetAllTransactions() ([]models.TransactionResponse, error) {
	transactions, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []models.TransactionResponse
	for _, t := range transactions {
		response = append(response, models.TransactionResponse{
			ID:            t.ID.Hex(),
			ProductID:     t.ProductID.Hex(),
			Date:          t.Date,
			Price:         t.Price,
			PaymentMethod: t.PaymentMethod,
			PaymentID:     t.PaymentID,
		})
	}

	if response == nil {
		response = []models.TransactionResponse{}
	}

	return response, nil
}

func (s *transactionService) GetTransactionByID(id string) (*models.TransactionResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	transaction, err := s.repo.FindByID(objectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	response := &models.TransactionResponse{
		ID:            transaction.ID.Hex(),
		ProductID:     transaction.ProductID.Hex(),
		Date:          transaction.Date,
		Price:         transaction.Price,
		PaymentMethod: transaction.PaymentMethod,
		PaymentID:     transaction.PaymentID,
	}

	return response, nil
}

func (s *transactionService) UpdateTransaction(id string, req *models.TransactionUpdateRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid transaction ID")
	}

	updateData := bson.M{}
	if req.ProductID != "" {
		productID, err := primitive.ObjectIDFromHex(req.ProductID)
		if err != nil {
			return errors.New("invalid product_id")
		}
		updateData["product_id"] = productID
	}
	if req.Price > 0 {
		updateData["price"] = req.Price
	}
	if req.PaymentMethod != "" {
		updateData["payment_method"] = req.PaymentMethod
	}

	if len(updateData) == 0 {
		return errors.New("no fields to update")
	}

	if err := s.repo.Update(objectID, updateData); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("transaction not found")
		}
		return err
	}

	return nil
}

func (s *transactionService) DeleteTransaction(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid transaction ID")
	}

	if err := s.repo.Delete(objectID); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("transaction not found")
		}
		return err
	}

	return nil
}
