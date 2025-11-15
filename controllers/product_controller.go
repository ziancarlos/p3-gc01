package controllers

import (
	"net/http"
	"p3-graded-challenge-1-ziancarlos/models"
	"p3-graded-challenge-1-ziancarlos/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ProductController struct {
	service   service.ProductService
	validator *validator.Validate
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{
		service:   service,
		validator: validator.New(),
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Add a new product to the database
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.ProductRequest true "Product data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (ctrl *ProductController) CreateProduct(c echo.Context) error {
	var req models.ProductRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := ctrl.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	response, err := ctrl.service.CreateProduct(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to create product",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Product created successfully",
		Data:    response,
	})
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve all products from the database
// @Tags products
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (ctrl *ProductController) GetAllProducts(c echo.Context) error {
	response, err := ctrl.service.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to retrieve products",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Products retrieved successfully",
		Data:    response,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Retrieve a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func (ctrl *ProductController) GetProductByID(c echo.Context) error {
	id := c.Param("id")

	response, err := ctrl.service.GetProductByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "Product not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product retrieved successfully",
		Data:    response,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.ProductRequest true "Product data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
func (ctrl *ProductController) UpdateProduct(c echo.Context) error {
	id := c.Param("id")

	var req models.ProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := ctrl.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	if err := ctrl.service.UpdateProduct(id, &req); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to update product",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product updated successfully",
	})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (ctrl *ProductController) DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	if err := ctrl.service.DeleteProduct(id); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to delete product",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product deleted successfully",
	})
}

