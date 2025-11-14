package service

import (
	"payment-service/models"
	"payment-service/repository"

	"github.com/go-playground/validator/v10"
)

type PaymentService interface {
	CreatePayment(req *models.PaymentRequest) (*models.PaymentResponse, error)
}

type paymentService struct {
	repo      repository.PaymentRepository
	validator *validator.Validate
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *paymentService) CreatePayment(req *models.PaymentRequest) (*models.PaymentResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	payment := &models.Payment{
		Amount: req.Amount,
	}
	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}
	return &models.PaymentResponse{
		ID:     payment.ID.Hex(),
		Amount: payment.Amount,
	}, nil
}
