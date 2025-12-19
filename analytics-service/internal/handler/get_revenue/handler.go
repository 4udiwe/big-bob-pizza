package get_revenue

import (
	"context"
	"net/http"
	"time"

	h "github.com/4udiwe/big-bob-pizza/analytics-service/internal/handler"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s AnalyticsService
}

func New(s AnalyticsService) h.Handler {
	return &handler{s: s}
}

type AnalyticsService interface {
	GetRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error)
}

// GetRevenue godoc
// @Summary Получить выручку за период
// @Description Возвращает общую выручку от созданных заказов за указанный период
// @Tags analytics
// @Accept json
// @Produce json
// @Param startDate query string true "Начальная дата (RFC3339)" format(date-time)
// @Param endDate query string true "Конечная дата (RFC3339)" format(date-time)
// @Success 200 {object} RevenueResponse
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /analytics/revenue [get]
func (h *handler) Handle(c echo.Context) error {
	startDateStr := c.QueryParam("startDate")
	endDateStr := c.QueryParam("endDate")

	if startDateStr == "" || endDateStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "startDate and endDate are required")
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid startDate format, expected RFC3339")
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid endDate format, expected RFC3339")
	}

	if endDate.Before(startDate) {
		return echo.NewHTTPError(http.StatusBadRequest, "endDate must be after startDate")
	}

	revenue, err := h.s.GetRevenue(c.Request().Context(), startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := RevenueResponse{
		Revenue:   revenue,
		StartDate: startDate,
		EndDate:   endDate,
	}

	return c.JSON(http.StatusOK, resp)
}

type RevenueResponse struct {
	Revenue   float64   `json:"revenue"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

