package get_orders_by_user

import (
	"net/http"
	"strconv"
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

// GetOrdersByUser godoc
// @Summary Получить заказы пользователя
// @Description Возвращает список заказов пользователя с пагинацией
// @Tags orders
// @Accept json
// @Produce json
// @Param userId path string true "ID пользователя (UUID)"
// @Param limit query int false "Количество записей на странице" default(20) minimum(1) maximum(100)
// @Param offset query int false "Смещение для пагинации" default(0) minimum(0)
// @Success 200 {object} OrdersResponse
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /orders/user/{userId} [get]
func (h *handler) Handle(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	// Парсим параметры пагинации
	limit := 20
	offset := 0

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid limit parameter")
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid offset parameter")
		}
	}

	orders, total, err := h.s.GetOrdersByUser(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := OrdersResponse{
		Orders: make([]OrderResponse, len(orders)),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	for i, order := range orders {
		resp.Orders[i] = OrderResponse{
			ID:          order.ID,
			CustomerID:  order.CustomerID,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			Currency:    order.Currency,
			PaymentID:   order.PaymentID,
			DeliveryID:  order.DeliveryID,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
			Items: lo.Map(order.Items, func(item entity.OrderItem, _ int) OrderItemResponse {
				return OrderItemResponse{
					ID:           item.ID,
					ProductID:    item.ProductID,
					ProductName:  item.ProductName,
					ProductPrice: item.ProductPrice,
					Amount:       item.Amount,
					TotalPrice:   item.TotalPrice,
					Notes:        item.Notes,
				}
			}),
		}
	}

	return c.JSON(http.StatusOK, resp)
}

type OrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

type OrderResponse struct {
	ID          uuid.UUID           `json:"id"`
	CustomerID  uuid.UUID           `json:"customerId"`
	Status      entity.OrderStatus  `json:"status"`
	TotalAmount float64             `json:"totalAmount"`
	Currency    string              `json:"currency"`
	PaymentID   *uuid.UUID          `json:"paymentId,omitempty"`
	DeliveryID  *uuid.UUID          `json:"deliveryId,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Items       []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"productId"`
	ProductName  string    `json:"productName"`
	ProductPrice float64   `json:"productPrice"`
	Amount       int       `json:"amount"`
	TotalPrice   float64   `json:"totalPrice"`
	Notes        string    `json:"notes"`
}
