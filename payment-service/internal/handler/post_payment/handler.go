package post_payment

import (
	"errors"
	"net/http"
	"time"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	h "github.com/4udiwe/big-bob-pizza/payment-service/internal/handler"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/handler/decorator"
	service "github.com/4udiwe/big-bob-pizza/payment-service/internal/service/payment"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s PaymentService
}

func New(s PaymentService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct {
	OrderID uuid.UUID `json:"orderId" validate:"required"`
	Amount  float64   `json:"amount" validate:"required,min=0"`
}

type Response struct {
	ID            uuid.UUID            `json:"id"`
	OrderID       uuid.UUID            `json:"orderId"`
	Amount        float64              `json:"amount"`
	Currency      string               `json:"currency"`
	Status        entity.PaymentStatus `json:"status"`
	FailureReason *string              `json:"failureReason,omitempty"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}

// ProcessPayment godoc
// @Summary Провести оплату заказа
// @Description Обрабатывает платеж для указанного заказа. Заказ должен быть доступен для оплаты (создан не более 30 минут назад)
// @Tags payments
// @Accept json
// @Produce json
// @Param request body Request true "Данные для оплаты"
// @Success 200 {object} Response
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /payments [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	payment, err := h.s.ProcessPayment(c.Request().Context(), in.OrderID, in.Amount)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrOrderAlreadyPaid) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if errors.Is(err, service.ErrInvalidAmount) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := Response{
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
