package payment_repository

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{Postgres: pg}
}

func (r *Repository) Create(ctx context.Context, payment entity.Payment) (entity.Payment, error) {
	logrus.Infof("PaymentRepository.Create: orderID=%s amount=%.2f", payment.OrderID, payment.Amount)

	query, args, _ := r.Builder.
		Insert("payments").
		Columns("order_id", "amount", "currency", "status_id").
		Values(payment.OrderID, payment.Amount, payment.Currency, squirrel.Expr("(SELECT id FROM payment_status WHERE name = ?)", payment.Status.Name)).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()

	row := r.GetTxManager(ctx).QueryRow(ctx, query, args...)
	if err := row.Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt); err != nil {
		logrus.Errorf("PaymentRepository.Create: scan error: %v", err)
		return entity.Payment{}, err
	}

	logrus.Infof("PaymentRepository.Create: created paymentID=%s", payment.ID)
	return payment, nil
}

func (r *Repository) GetByID(ctx context.Context, paymentID uuid.UUID) (entity.Payment, error) {
	query := `
		SELECT
			p.id, p.order_id, p.amount, p.currency,
			p.status_id, s.name AS status_name,
			p.failure_reason, p.created_at, p.updated_at
		FROM payments p
		JOIN payment_status s ON s.id = p.status_id
		WHERE p.id = $1
	`

	row := r.GetTxManager(ctx).QueryRow(ctx, query, paymentID)
	var dto RowPayment
	if err := row.Scan(
		&dto.ID, &dto.OrderID, &dto.Amount, &dto.Currency,
		&dto.StatusID, &dto.StatusName,
		&dto.FailureReason, &dto.CreatedAt, &dto.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Payment{}, ErrPaymentNotFound
		}
		return entity.Payment{}, err
	}

	return dto.ToEntity(), nil
}

func (r *Repository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (entity.Payment, error) {
	query := `
		SELECT
			p.id, p.order_id, p.amount, p.currency,
			p.status_id, s.name AS status_name,
			p.failure_reason, p.created_at, p.updated_at
		FROM payments p
		JOIN payment_status s ON s.id = p.status_id
		WHERE p.order_id = $1
		ORDER BY p.created_at DESC
		LIMIT 1
	`

	row := r.GetTxManager(ctx).QueryRow(ctx, query, orderID)
	var dto RowPayment
	if err := row.Scan(
		&dto.ID, &dto.OrderID, &dto.Amount, &dto.Currency,
		&dto.StatusID, &dto.StatusName,
		&dto.FailureReason, &dto.CreatedAt, &dto.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Payment{}, ErrPaymentNotFound
		}
		return entity.Payment{}, err
	}

	return dto.ToEntity(), nil
}

func (r *Repository) UpdateStatus(ctx context.Context, paymentID uuid.UUID, status entity.PaymentStatus, failureReason *string) error {
	logrus.Infof("PaymentRepository.UpdateStatus: paymentID=%s status=%s", paymentID, status.Name)

	query, args, _ := r.Builder.
		Update("payments").
		Set("status_id", squirrel.Expr("(SELECT id FROM payment_status WHERE name = ?)", status.Name)).
		Set("updated_at", time.Now()).
		Where("id = ?", paymentID).
		ToSql()

	if failureReason != nil {
		query, args, _ = r.Builder.
			Update("payments").
			Set("status_id", squirrel.Expr("(SELECT id FROM payment_status WHERE name = ?)", status.Name)).
			Set("failure_reason", *failureReason).
			Set("updated_at", time.Now()).
			Where("id = ?", paymentID).
			ToSql()
	}

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("PaymentRepository.UpdateStatus: update error: %v", err)
		return err
	}

	return nil
}

// GetAllPayments возвращает список платежей с пагинацией и фильтрацией
func (r *Repository) GetAllPayments(ctx context.Context, limit, offset int, status *entity.PaymentStatusName, userID *uuid.UUID) ([]entity.PaymentWithUser, int, error) {
	logrus.Infof("PaymentRepository.GetAllPayments: limit=%d offset=%d", limit, offset)

	// Строим условия для WHERE
	whereConditions := squirrel.And{}
	if status != nil {
		whereConditions = append(whereConditions, squirrel.Expr("s.name = ?", string(*status)))
	}
	if userID != nil {
		whereConditions = append(whereConditions, squirrel.Eq{"oc.user_id": *userID})
	}

	// Получаем общее количество
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("payments p").
		Join("payment_status s ON s.id = p.status_id").
		LeftJoin("order_cache oc ON oc.order_id = p.order_id")

	if len(whereConditions) > 0 {
		countQuery = countQuery.Where(whereConditions)
	}

	countSQL, countArgs, _ := countQuery.ToSql()

	var total int
	err := r.GetTxManager(ctx).QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		logrus.Errorf("PaymentRepository.GetAllPayments: count error: %v", err)
		return nil, 0, err
	}

	// Получаем страницу данных
	query := r.Builder.
		Select(`
			p.id, p.order_id, p.amount, p.currency,
			p.status_id, s.name AS status_name,
			p.failure_reason, p.created_at, p.updated_at,
			oc.user_id
		`).
		From("payments p").
		Join("payment_status s ON s.id = p.status_id").
		LeftJoin("order_cache oc ON oc.order_id = p.order_id")

	if len(whereConditions) > 0 {
		query = query.Where(whereConditions)
	}

	querySQL, args, _ := query.
		OrderBy("p.created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, querySQL, args...)
	if err != nil {
		logrus.Errorf("PaymentRepository.GetAllPayments: query error: %v", err)
		return nil, 0, err
	}

	payments, err := pgx.CollectRows(rows, pgx.RowToStructByName[RowPaymentWithUser])
	if err != nil {
		logrus.Errorf("PaymentRepository.GetAllPayments: scan error: %v", err)
		return nil, 0, err
	}

	result := make([]entity.PaymentWithUser, len(payments))
	for i, p := range payments {
		result[i] = p.ToPaymentWithUser()
	}

	logrus.Infof("PaymentRepository.GetAllPayments: fetched %d payments", len(result))
	return result, total, nil
}

// GetPaymentsByUserID возвращает платежи пользователя с пагинацией
func (r *Repository) GetPaymentsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.PaymentWithUser, int, error) {
	return r.GetAllPayments(ctx, limit, offset, nil, &userID)
}
