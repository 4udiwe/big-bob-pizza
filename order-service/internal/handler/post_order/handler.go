package post_order

import (
	"errors"
	"net/http"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	h "github.com/4udiwe/big-bob-pizza/order-service/internal/handler"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/handler/decorator"
	service "github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s OrderService
}

func New(s OrderService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct {
	CustomerID  uuid.UUID          `json:"customerId" validate:"required"`
	TotalAmount float64            `json:"totalAmount" validate:"required,min=0"`
	Currency    string             `json:"currency" validate:"required"`
	Items       []RequestOrderItem `json:"items" validate:"required"`
}

type RequestOrderItem struct {
	ProductID    uuid.UUID `json:"productId" validate:"required"`
	ProductName  string    `json:"productName" validate:"required"`
	ProductPrice float64   `json:"productPrice" validate:"required,min=0"`
	Amount       int       `json:"amount" validate:"required,min=1"`
	TotalPrice   float64   `json:"totalPrice" validate:"required,min=0"`
	Notes        string    `json:"notes"`
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

// CreateOrder godoc
// @Summary Создать новый заказ
// @Description Создает новый заказ для указанного пользователя. После создания публикуется событие order.created
// @Tags orders
// @Accept json
// @Produce json
// @Param request body Request true "Данные заказа"
// @Success 201 {object} Response
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 409 {string} string "Заказ уже существует"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /orders [post]
func (h *handler) Handle(c echo.Context, in Request) error {

	order := entity.Order{
		CustomerID:  in.CustomerID,
		TotalAmount: in.TotalAmount,
		Currency:    in.Currency,
		Items: lo.Map(in.Items, func(i RequestOrderItem, _ int) entity.OrderItem {
			return entity.OrderItem{
				ProductID:    i.ProductID,
				ProductName:  i.ProductName,
				ProductPrice: i.ProductPrice,
				Amount:       i.Amount,
				TotalPrice:   i.TotalPrice,
				Notes:        i.Notes,
			}
		}),
	}

	offer, err := h.s.CreateOrder(c.Request().Context(), order)
	if err != nil {
		if errors.Is(err, service.ErrOrderAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := Response{
		ID:          offer.ID,
		CustomerID:  offer.CustomerID,
		Status:      offer.Status,
		TotalAmount: offer.TotalAmount,
		Currency:    offer.Currency,
		PaymentID:   offer.PaymentID,
		DeliveryID:  offer.DeliveryID,
		CreatedAt:   offer.CreatedAt,
		UpdatedAt:   offer.UpdatedAt,
		Items: lo.Map(offer.Items, func(i entity.OrderItem, _ int) ResponseOrderItem {
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

	return c.JSON(http.StatusCreated, resp)
}
