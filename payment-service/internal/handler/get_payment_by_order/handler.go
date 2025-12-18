package get_payment_by_order

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	h "github.com/4udiwe/big-bob-pizza/payment-service/internal/handler"
	service "github.com/4udiwe/big-bob-pizza/payment-service/internal/service/payment"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s PaymentService
}

func New(s PaymentService) h.Handler {
	return &handler{s: s}
}

type PaymentService interface {
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (entity.Payment, error)
}

// GetPaymentByOrder godoc
// @Summary Получить платеж по ID заказа
// @Description Возвращает информацию о платеже для указанного заказа
// @Tags payments
// @Accept json
// @Produce json
// @Param orderId path string true "ID заказа (UUID)"
// @Success 200 {object} PaymentResponse
// @Failure 404 {string} string "Платеж не найден"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /payments/order/{orderId} [get]
func (h *handler) Handle(c echo.Context) error {
	orderIDStr := c.Param("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order ID")
	}

	payment, err := h.s.GetPaymentByOrderID(c.Request().Context(), orderID)
	if err != nil {
		if errors.Is(err, service.ErrPaymentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        payment.Status,
		FailureReason: payment.FailureReason,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}

	return c.JSON(http.StatusOK, resp)
}

type PaymentResponse struct {
	ID            uuid.UUID            `json:"id"`
	OrderID       uuid.UUID            `json:"orderId"`
	Amount        float64              `json:"amount"`
	Currency      string               `json:"currency"`
	Status        entity.PaymentStatus `json:"status"`
	FailureReason *string              `json:"failureReason,omitempty"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}
