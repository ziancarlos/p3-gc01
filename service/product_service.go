package service

import (
	"errors"
	"p3-graded-challenge-1-ziancarlos/models"
	"p3-graded-challenge-1-ziancarlos/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService interface {
	CreateProduct(req *models.ProductRequest) (*models.ProductResponse, error)
	GetAllProducts() ([]models.ProductResponse, error)
	GetProductByID(id string) (*models.ProductResponse, error)
	UpdateProduct(id string, req *models.ProductRequest) error
	DeleteProduct(id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(req *models.ProductRequest) (*models.ProductResponse, error) {
	product := &models.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	response := &models.ProductResponse{
		ID:    product.ID.Hex(),
		Name:  product.Name,
		Price: product.Price,
	}

	return response, nil
}

func (s *productService) GetAllProducts() ([]models.ProductResponse, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []models.ProductResponse
	for _, p := range products {
		response = append(response, models.ProductResponse{
			ID:    p.ID.Hex(),
			Name:  p.Name,
			Price: p.Price,
		})
	}

	if response == nil {
		response = []models.ProductResponse{}
	}

	return response, nil
}

func (s *productService) GetProductByID(id string) (*models.ProductResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	product, err := s.repo.FindByID(objectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	response := &models.ProductResponse{
		ID:    product.ID.Hex(),
		Name:  product.Name,
		Price: product.Price,
	}

	return response, nil
}

func (s *productService) UpdateProduct(id string, req *models.ProductRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	if err := s.repo.Update(objectID, req); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("product not found")
		}
		return err
	}

	return nil
}

func (s *productService) DeleteProduct(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	if err := s.repo.Delete(objectID); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("product not found")
		}
		return err
	}

	return nil
}
