package controllers

import (
	"net/http"
	"payment-service/models"
	"payment-service/service"

	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	service service.PaymentService
}

func NewPaymentController(service service.PaymentService) *PaymentController {
	return &PaymentController{service: service}
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Add a new payment to the database
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body models.PaymentRequest true "Payment data (amount)"
// @Success 201 {object} models.PaymentResponse
// @Failure 400 {object} map[string]string
// @Router /payments [post]
func (ctrl *PaymentController) CreatePayment(c echo.Context) error {
	var req models.PaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := ctrl.service.CreatePayment(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, resp)
}
