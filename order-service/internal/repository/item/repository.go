package item_repository

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{Postgres: postgres}
}

// Inserts items for order.
// Returns filled items (with TotalPrice and ID).
func (r *Repository) InsertItems(ctx context.Context, orderID uuid.UUID, items []entity.OrderItem) ([]entity.OrderItem, error) {
	logrus.Infof("ItemRepository.InsertItems: insert for orderID=%v", orderID)

	builder := r.Builder.
		Insert("order_item").
		Columns("order_id", "product_id", "product_name", "product_price", "amount", "notes")

	for _, i := range items {
		builder = builder.Values(orderID, i.ProductID, i.ProductName, i.ProductPrice, i.Amount, i.Notes)
	}

	query, args, _ := builder.
		Suffix(`RETURNING 
			id,
			order_id,
			product_id,
			product_name,
			product_price,
			amount,
			total_price,
			notes`).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Infof("ItemRepository.InsertItems: query error: %v", err)
		return nil, err
	}

	rowsItem, err := pgx.CollectRows(rows, pgx.RowToStructByName[RowItem])
	if err != nil {
		logrus.Infof("ItemRepository.InsertItems: scan error: %v", err)
		return nil, err
	}

	entities := lo.Map(rowsItem, func(r RowItem, _ int) entity.OrderItem { return r.ToEntity() })

	logrus.Infof("ItemRepository.InsertItems: items inserted for orderID=%v", orderID)
	return entities, nil
}
