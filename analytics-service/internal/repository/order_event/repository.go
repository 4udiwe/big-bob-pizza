package order_event

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{Postgres: pg}
}

// Save сохраняет событие заказа. Если событие с таким event_id уже существует, возвращает nil (идемпотентность)
func (r *Repository) Save(ctx context.Context, event entity.OrderEvent) error {
	logrus.Infof("OrderEventRepository.Save: eventType=%s orderID=%s eventID=%s", event.EventType, event.OrderID, event.EventID)

	var dto OrderEventDTO
	dto.FromEntity(event)
	dto.ID = uuid.New()

	// Используем raw SQL для правильной обработки NULL значений
	query := `
		INSERT INTO order_events (id, event_id, event_type, order_id, user_id, amount, payment_id, reason, occurred_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var userID, paymentID interface{}
	var amount interface{}
	var reason interface{}

	if dto.UserID.Valid {
		userID = dto.UserID.String
	} else {
		userID = nil
	}

	if dto.Amount.Valid {
		amount = dto.Amount.Float64
	} else {
		amount = nil
	}

	if dto.PaymentID.Valid {
		paymentID = dto.PaymentID.String
	} else {
		paymentID = nil
	}

	if dto.Reason.Valid {
		reason = dto.Reason.String
	} else {
		reason = nil
	}

	args := []interface{}{dto.ID, dto.EventID, dto.EventType, dto.OrderID, userID, amount, paymentID, reason, dto.OccurredAt, time.Now()}

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			// Событие уже существует (дедупликация), это нормально
			logrus.Debugf("OrderEventRepository.Save: event already exists eventID=%s", event.EventID)
			return nil
		}
		logrus.Errorf("OrderEventRepository.Save: insert error: %v", err)
		return err
	}

	logrus.Infof("OrderEventRepository.Save: saved eventID=%s", event.EventID)
	return nil
}

// GetByOrderID возвращает все события для заказа
func (r *Repository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.OrderEvent, error) {
	query := `
		SELECT id, event_id, event_type, order_id, user_id, amount, payment_id, reason, occurred_at, created_at
		FROM order_events
		WHERE order_id = $1
		ORDER BY occurred_at ASC
	`

	rows, err := r.GetTxManager(ctx).Query(ctx, query, orderID)
	if err != nil {
		logrus.Errorf("OrderEventRepository.GetByOrderID: query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []entity.OrderEvent
	for rows.Next() {
		var dto OrderEventDTO
		err := rows.Scan(
			&dto.ID,
			&dto.EventID,
			&dto.EventType,
			&dto.OrderID,
			&dto.UserID,
			&dto.Amount,
			&dto.PaymentID,
			&dto.Reason,
			&dto.OccurredAt,
			&dto.CreatedAt,
		)
		if err != nil {
			logrus.Errorf("OrderEventRepository.GetByOrderID: scan error: %v", err)
			return nil, err
		}

		events = append(events, dto.ToEntity())
	}

	return events, nil
}

// GetStatsByDateRange возвращает статистику за период
func (r *Repository) GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]OrderStats, error) {
	query := `
		SELECT 
			DATE_TRUNC('day', occurred_at) AS date,
			event_type,
			COUNT(*) AS count,
			COUNT(DISTINCT order_id) AS unique_orders,
			COUNT(DISTINCT user_id) AS unique_users,
			SUM(amount) AS total_amount
		FROM order_events
		WHERE occurred_at >= $1 AND occurred_at < $2
		GROUP BY DATE_TRUNC('day', occurred_at), event_type
		ORDER BY date DESC, event_type
	`

	rows, err := r.GetTxManager(ctx).Query(ctx, query, startDate, endDate)
	if err != nil {
		logrus.Errorf("OrderEventRepository.GetStatsByDateRange: query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var stats []OrderStats
	for rows.Next() {
		var s OrderStats
		var totalAmount sql.NullFloat64
		err := rows.Scan(
			&s.Date,
			&s.EventType,
			&s.Count,
			&s.UniqueOrders,
			&s.UniqueUsers,
			&totalAmount,
		)
		if err != nil {
			logrus.Errorf("OrderEventRepository.GetStatsByDateRange: scan error: %v", err)
			return nil, err
		}

		if totalAmount.Valid {
			s.TotalAmount = &totalAmount.Float64
		}

		stats = append(stats, s)
	}

	return stats, nil
}

// GetTotalRevenue возвращает общую выручку за период
func (r *Repository) GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM order_events
		WHERE event_type = 'order.created'
			AND occurred_at >= $1 AND occurred_at < $2
			AND amount IS NOT NULL
	`

	var total float64
	err := r.GetTxManager(ctx).QueryRow(ctx, query, startDate, endDate).Scan(&total)
	if err != nil {
		logrus.Errorf("OrderEventRepository.GetTotalRevenue: query error: %v", err)
		return 0, err
	}

	return total, nil
}

type OrderStats struct {
	Date         time.Time
	EventType    string
	Count        int
	UniqueOrders int
	UniqueUsers  int
	TotalAmount  *float64
}

