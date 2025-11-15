package controllers

import (
	"net/http"
	"p3-graded-challenge-1-ziancarlos/models"
	"p3-graded-challenge-1-ziancarlos/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	service   service.TransactionService
	validator *validator.Validate
}

func NewTransactionController(service service.TransactionService) *TransactionController {
	return &TransactionController{
		service:   service,
		validator: validator.New(),
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Add a new transaction to the database
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.TransactionRequest true "Transaction data (must include product_id)"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions [post]
func (ctrl *TransactionController) CreateTransaction(c echo.Context) error {
	var req models.TransactionRequest

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

	response, err := ctrl.service.CreateTransaction(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to create transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Transaction created successfully",
		Data:    response,
	})
}

// GetAllTransactions godoc
// @Summary Get all transactions
// @Description Retrieve all transactions from the database
// @Tags transactions
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions [get]
func (ctrl *TransactionController) GetAllTransactions(c echo.Context) error {
	response, err := ctrl.service.GetAllTransactions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to retrieve transactions",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Transactions retrieved successfully",
		Data:    response,
	})
}

// GetTransactionByID godoc
// @Summary Get transaction by ID
// @Description Retrieve a transaction by its ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /transactions/{id} [get]
func (ctrl *TransactionController) GetTransactionByID(c echo.Context) error {
	id := c.Param("id")

	response, err := ctrl.service.GetTransactionByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Transaction retrieved successfully",
		Data:    response,
	})
}

// UpdateTransaction godoc
// @Summary Update a transaction
// @Description Update an existing transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param transaction body models.TransactionUpdateRequest true "Transaction data (can include product_id)"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions/{id} [put]
func (ctrl *TransactionController) UpdateTransaction(c echo.Context) error {
	id := c.Param("id")

	var req models.TransactionUpdateRequest
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

	if err := ctrl.service.UpdateTransaction(id, &req); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to update transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Transaction updated successfully",
	})
}

// DeleteTransaction godoc
// @Summary Delete a transaction
// @Description Delete a transaction by its ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions/{id} [delete]
func (ctrl *TransactionController) DeleteTransaction(c echo.Context) error {
	id := c.Param("id")

	if err := ctrl.service.DeleteTransaction(id); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to delete transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Transaction deleted successfully",
	})
}

