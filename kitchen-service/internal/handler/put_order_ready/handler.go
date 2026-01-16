package put_order_ready

import (
	"context"
	"net/http"

	h "github.com/4udiwe/big-bob-pizza/kitchen-service/internal/handler"
	"github.com/4udiwe/big-bob-pizza/kitchen-service/internal/handler/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderService interface {
	MarkOrderReady(ctx context.Context, orderID uuid.UUID) error
}

type handler struct {
	s OrderService
}

func New(s OrderService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct {
	OrderID uuid.UUID `json:"orderId" validate:"required"`
}

// MarkOrderReady godoc
// @Summary Отметить заказ как готовый
// @Description Отмечает заказ как готовый для доставки
// @Tags orders
// @Accept json
// @Produce json
// @Param request body Request true "Данные заказа"
// @Success 202
// @Failure 500 {object} echo.HTTPError "Внутренняя ошибка сервиса"
// @Router /kitchen [put]
func (h *handler) Handle(c echo.Context, in Request) error {
	err := h.s.MarkOrderReady(c.Request().Context(), in.OrderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusAccepted)
}
