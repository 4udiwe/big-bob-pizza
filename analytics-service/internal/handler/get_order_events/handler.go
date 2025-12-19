package get_order_events

import (
	"context"
	"net/http"
	"time"

	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/entity"
	h "github.com/4udiwe/big-bob-pizza/analytics-service/internal/handler"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s AnalyticsService
}

func New(s AnalyticsService) h.Handler {
	return &handler{s: s}
}

type AnalyticsService interface {
	GetOrderEvents(ctx context.Context, orderID uuid.UUID) ([]entity.OrderEvent, error)
}

// GetOrderEvents godoc
// @Summary Получить события заказа
// @Description Возвращает все события для указанного заказа
// @Tags analytics
// @Accept json
// @Produce json
// @Param orderId path string true "ID заказа (UUID)" format(uuid)
// @Success 200 {object} OrderEventsResponse
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /analytics/orders/{orderId}/events [get]
func (h *handler) Handle(c echo.Context) error {
	orderIDStr := c.Param("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid orderId format")
	}

	events, err := h.s.GetOrderEvents(c.Request().Context(), orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := OrderEventsResponse{
		OrderID: orderID,
		Events:  make([]OrderEventResponse, len(events)),
	}

	for i, e := range events {
		resp.Events[i] = OrderEventResponse{
			ID:         e.ID,
			EventID:    e.EventID,
			EventType:  e.EventType,
			OrderID:    e.OrderID,
			UserID:     e.UserID,
			Amount:     e.Amount,
			PaymentID:  e.PaymentID,
			Reason:     e.Reason,
			OccurredAt: e.OccurredAt,
			CreatedAt:  e.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

type OrderEventsResponse struct {
	OrderID uuid.UUID            `json:"orderId"`
	Events  []OrderEventResponse `json:"events"`
}

type OrderEventResponse struct {
	ID         uuid.UUID  `json:"id"`
	EventID    uuid.UUID  `json:"eventId"`
	EventType  string     `json:"eventType"`
	OrderID    uuid.UUID  `json:"orderId"`
	UserID     *uuid.UUID `json:"userId,omitempty"`
	Amount     *float64   `json:"amount,omitempty"`
	PaymentID  *uuid.UUID `json:"paymentId,omitempty"`
	Reason     *string    `json:"reason,omitempty"`
	OccurredAt time.Time  `json:"occurredAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

