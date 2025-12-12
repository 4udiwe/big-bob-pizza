package get_payments

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	h "github.com/4udiwe/big-bob-pizza/payment-service/internal/handler"
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
	GetAllPayments(ctx context.Context, limit, offset int, status *entity.PaymentStatusName, userID *uuid.UUID) ([]entity.PaymentWithUser, int, error)
}

// GetPayments godoc
// @Summary Получить список всех платежей
// @Description Возвращает список платежей с пагинацией. Можно фильтровать по статусу и пользователю
// @Tags payments
// @Accept json
// @Produce json
// @Param limit query int false "Количество записей на странице" default(20) minimum(1) maximum(100)
// @Param offset query int false "Смещение для пагинации" default(0) minimum(0)
// @Param status query string false "Фильтр по статусу: pending, completed, failed"
// @Param userId query string false "Фильтр по ID пользователя (UUID)"
// @Success 200 {object} PaymentsResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /payments [get]
func (h *handler) Handle(c echo.Context) error {
	// Парсим параметры пагинации
	limit := 20
	offset := 0
	var err error

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

	// Парсим фильтры
	var status *entity.PaymentStatusName
	if statusStr := c.QueryParam("status"); statusStr != "" {
		s := entity.PaymentStatusName(statusStr)
		if s != entity.PaymentStatusPending && s != entity.PaymentStatusCompleted && s != entity.PaymentStatusFailed {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid status parameter")
		}
		status = &s
	}

	var userID *uuid.UUID
	if userIDStr := c.QueryParam("userId"); userIDStr != "" {
		id, err := uuid.Parse(userIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid userId parameter")
		}
		userID = &id
	}

	payments, total, err := h.s.GetAllPayments(c.Request().Context(), limit, offset, status, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := PaymentsResponse{
		Payments: make([]PaymentResponse, len(payments)),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	for i, p := range payments {
		resp.Payments[i] = PaymentResponse{
			ID:            p.Payment.ID,
			OrderID:       p.Payment.OrderID,
			Amount:        p.Payment.Amount,
			Currency:      p.Payment.Currency,
			Status:        p.Payment.Status,
			FailureReason: p.Payment.FailureReason,
			UserID:        p.UserID,
			CreatedAt:     p.Payment.CreatedAt,
			UpdatedAt:     p.Payment.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

type PaymentsResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type PaymentResponse struct {
	ID            uuid.UUID            `json:"id"`
	OrderID       uuid.UUID            `json:"orderId"`
	Amount        float64              `json:"amount"`
	Currency      string               `json:"currency"`
	Status        entity.PaymentStatus `json:"status"`
	FailureReason *string              `json:"failureReason,omitempty"`
	UserID        *uuid.UUID           `json:"userId,omitempty"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}
