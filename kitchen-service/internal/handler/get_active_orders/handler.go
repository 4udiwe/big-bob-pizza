package get_active_orders

import (
	"context"
	"net/http"

	"github.com/4udiwe/big-bob-pizza/kitchen-service/internal/entity"
	h "github.com/4udiwe/big-bob-pizza/kitchen-service/internal/handler"
	"github.com/4udiwe/big-bob-pizza/kitchen-service/internal/handler/decorator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type OrderService interface {
	GetActiveOrders(ctx context.Context) ([]entity.KitchenOrder, error)
}

type handler struct {
	s OrderService
}

func New(s OrderService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct{}

type Response struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID      string `json:"id"`
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
	Items   []Item `json:"items"`
}

type Item struct {
	ItemID   string `json:"itemId"`
	ItemName string `json:"itemName"`
	Quantity int    `json:"quantity"`
}

// GetActiveOrders godoc
// @Summary Получить активные заказы
// @Description Получает список активных заказов на кухне
// @Tags orders
// @Accept json
// @Produce json
// @Param request body Request true "Данные заказа"
// @Success 200 {object} Response
// @Failure 500 {object} echo.HTTPError "Внутренняя ошибка сервиса"
// @Router /kitchen [put]
func (h *handler) Handle(c echo.Context, in Request) error {
	orders, err := h.s.GetActiveOrders(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := Response{
		Orders: make([]Order, len(orders)),
	}

	for i, order := range orders {
		response.Orders[i] = Order{
			ID:      order.ID,
			OrderID: order.OrderID,
			Status:  string(order.Status),
			Items: lo.Map(order.Items, func(item entity.KitchenItem, _ int) Item {
				return Item{
					ItemID:   item.DishID,
					ItemName: item.Name,
					Quantity: item.Quantity,
				}
			}),
		}
	}

	return c.JSON(http.StatusOK, response)
}
