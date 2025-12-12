package get_order

import (
	"net/http"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	h "github.com/4udiwe/big-bob-pizza/order-service/internal/handler"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s OrderService
}

func New(s OrderService) h.Handler {
	return &handler{s: s}
}

// GetOrder godoc
// @Summary Получить заказ по ID
// @Description Возвращает информацию о заказе по его идентификатору
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "ID заказа (UUID)"
// @Success 200 {object} Response
// @Failure 404 {string} string "Заказ не найден"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /orders/{id} [get]
func (h *handler) Handle(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order ID")
	}

	order, err := h.s.GetOrderByID(c.Request().Context(), orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := Response{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Currency:    order.Currency,
		PaymentID:   order.PaymentID,
		DeliveryID:  order.DeliveryID,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Items: lo.Map(order.Items, func(i entity.OrderItem, _ int) ResponseOrderItem {
			return ResponseOrderItem{
				ID:           i.ID,
				ProductID:    i.ProductID,
				ProductName:  i.ProductName,
				ProductPrice: i.ProductPrice,
				Amount:       i.Amount,
				TotalPrice:   i.TotalPrice,
				Notes:        i.Notes,
			}
		}),
	}

	return c.JSON(http.StatusOK, resp)
}

type Response struct {
	ID          uuid.UUID           `json:"id"`
	CustomerID  uuid.UUID           `json:"customerId"`
	Status      entity.OrderStatus  `json:"status"`
	TotalAmount float64             `json:"totalAmount"`
	Currency    string              `json:"currency"`
	PaymentID   *uuid.UUID          `json:"paymentId,omitempty"`
	DeliveryID  *uuid.UUID          `json:"deliveryId,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Items       []ResponseOrderItem `json:"items"`
}

type ResponseOrderItem struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"productId"`
	ProductName  string    `json:"productName"`
	ProductPrice float64   `json:"productPrice"`
	Amount       int       `json:"amount"`
	TotalPrice   float64   `json:"totalPrice"`
	Notes        string    `json:"notes"`
}
