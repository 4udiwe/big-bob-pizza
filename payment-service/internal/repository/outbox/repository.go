package outbox_repository

import (
	"context"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{Postgres: pg}
}

func (r *Repository) Create(ctx context.Context, ev entity.OutboxEvent) error {
	logrus.Infof("OutboxRepository.Create: aggregate=%s id=%s type=%s",
		ev.AggregateType, ev.AggregateID, ev.EventType)

	query, args, _ := r.Builder.
		Insert("outbox").
		Columns("aggregate_type", "aggregate_id", "event_type", "payload").
		Values(ev.AggregateType, ev.AggregateID, ev.EventType, ev.Payload).
		Suffix("RETURNING id").
		ToSql()

	row := r.GetTxManager(ctx).QueryRow(ctx, query, args...)
	if err := row.Scan(&ev.ID); err != nil {
		logrus.Errorf("OutboxRepository.Create: scan error: %v", err)
		return err
	}

	logrus.Infof("OutboxRepository.Create: created eventID=%s", ev.ID)
	return nil
}

func (r *Repository) FetchPending(ctx context.Context, limit int) ([]outbox.Event, error) {
	logrus.Infof("OutboxRepository.FetchPending: limit=%d", limit)

	query := `
		SELECT
			o.id, o.aggregate_type, o.aggregate_id, o.event_type, o.payload,
			o.status_id,
			s.name AS status_name,
			o.created_at, o.processed_at
		FROM outbox o
		JOIN outbox_status s ON s.id = o.status_id
		WHERE s.name = $1
		ORDER BY o.created_at
		LIMIT $2
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.GetTxManager(ctx).Query(ctx, query, entity.OutboxStatusPending, limit)
	if err != nil {
		logrus.Errorf("OutboxRepository.FetchPending: query error: %v", err)
		return nil, err
	}

	dtoRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[RowOutbox])
	if err != nil {
		logrus.Errorf("OutboxRepository.FetchPending: scan error: %v", err)
		return nil, err
	}

	events := lo.Map(dtoRows, func(r RowOutbox, _ int) outbox.Event { return r.ToEvent() })

	logrus.Infof("OutboxRepository.FetchPending: fetched=%d", len(events))

	return events, nil
}

func (r *Repository) MarkProcessed(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	logrus.Infof("OutboxRepository.MarkProcessed: count=%d", len(ids))

	query, args, _ := r.Builder.
		Update("outbox").
		Set("status_id", squirrel.Expr("(SELECT id FROM outbox_status WHERE name = ?)", entity.OutboxStatusProcessed)).
		Set("processed_at", time.Now()).
		Where(squirrel.Eq{"id": ids}).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OutboxRepository.MarkProcessed: update error: %v", err)
		return err
	}

	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID, errorText string) error {
	logrus.Warnf("OutboxRepository.MarkFailed: id=%s err=%s", id, errorText)

	query, args, _ := r.Builder.
		Update("outbox").
		Set("status_id", squirrel.Expr("(SELECT id FROM outbox_status WHERE name = ?)", entity.OutboxStatusFailed)).
		Set("processed_at", time.Now()).
		Where("id = ?", id).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OutboxRepository.MarkFailed: update error: %v", err)
		return err
	}

	return nil
}

func (r *Repository) RequeueFailed(ctx context.Context, limit int) ([]outbox.Event, error) {
	logrus.Infof("OutboxRepository.RequeueFailed: limit=%d", limit)

	query := `
		SELECT
			o.id, o.aggregate_type, o.aggregate_id, o.event_type, o.payload,
			o.status_id,
			s.name AS status_name,
			o.created_at, o.processed_at
		FROM outbox o
		JOIN outbox_status s ON s.id = o.status_id
		WHERE s.name = $1
		ORDER BY o.created_at
		LIMIT $2
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.GetTxManager(ctx).Query(ctx, query, entity.OutboxStatusFailed, limit)
	if err != nil {
		logrus.Errorf("OutboxRepository.RequeueFailed: query error: %v", err)
		return nil, err
	}

	dtoRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[RowOutbox])
	if err != nil {
		logrus.Errorf("OutboxRepository.RequeueFailed: scan error: %v", err)
		return nil, err
	}

	events := lo.Map(dtoRows, func(r RowOutbox, _ int) outbox.Event { return r.ToEvent() })

	// Update status
	ids := lo.Map(events, func(e outbox.Event, _ int) uuid.UUID { return e.ID })

	query, args, _ := r.Builder.
		Update("outbox").
		Set("status_id", squirrel.Expr("(SELECT id FROM outbox_status WHERE name = ?)", entity.OutboxStatusPending)).
		Set("processed_at", nil).
		Where(squirrel.Expr("id = ANY(?)", ids)).
		ToSql()

	_, err = r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OutboxRepository.RequeueFailed: update error: %v", err)
		return nil, err
	}

	logrus.Infof("OutboxRepository.RequeueFailed: requeued=%d", len(events))
	return events, nil
}
