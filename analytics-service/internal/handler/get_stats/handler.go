package get_stats

import (
	"context"
	"net/http"
	"time"

	h "github.com/4udiwe/big-bob-pizza/analytics-service/internal/handler"
	order_event_repo "github.com/4udiwe/big-bob-pizza/analytics-service/internal/repository/order_event"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s AnalyticsService
}

func New(s AnalyticsService) h.Handler {
	return &handler{s: s}
}

type AnalyticsService interface {
	GetStats(ctx context.Context, startDate, endDate time.Time) ([]order_event_repo.OrderStats, error)
}

// GetStats godoc
// @Summary Получить статистику по заказам
// @Description Возвращает статистику по событиям заказов за указанный период
// @Tags analytics
// @Accept json
// @Produce json
// @Param startDate query string true "Начальная дата (RFC3339)" format(date-time)
// @Param endDate query string true "Конечная дата (RFC3339)" format(date-time)
// @Success 200 {object} StatsResponse
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /analytics/stats [get]
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

	stats, err := h.s.GetStats(c.Request().Context(), startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := StatsResponse{
		Stats:     make([]StatResponse, len(stats)),
		StartDate: startDate,
		EndDate:   endDate,
	}

	for i, s := range stats {
		resp.Stats[i] = StatResponse{
			Date:         s.Date,
			EventType:    s.EventType,
			Count:        s.Count,
			UniqueOrders: s.UniqueOrders,
			UniqueUsers:  s.UniqueUsers,
			TotalAmount:  s.TotalAmount,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

type StatsResponse struct {
	Stats     []StatResponse `json:"stats"`
	StartDate time.Time      `json:"startDate"`
	EndDate   time.Time      `json:"endDate"`
}

type StatResponse struct {
	Date         time.Time `json:"date"`
	EventType    string    `json:"eventType"`
	Count        int       `json:"count"`
	UniqueOrders int       `json:"uniqueOrders"`
	UniqueUsers  int       `json:"uniqueUsers"`
	TotalAmount  *float64  `json:"totalAmount,omitempty"`
}

