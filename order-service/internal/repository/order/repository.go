package order_repository

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/repository"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{Postgres: postgres}
}

// Inserts only order data (without order items).
// Receives entity with order data.
func (r *Repository) Create(ctx context.Context, order entity.Order) (entity.Order, error) {
	logrus.Infof("OrderRepository.Create: customerID=%v", order.CustomerID)

	query, args, _ := r.Builder.
		Insert("orders").
		Columns("customer_id", "total_amount", "currency").
		Values(order.CustomerID, order.TotalAmount, order.Currency).
		Suffix(`RETURNING 
				id,
				customer_id,
				status_id,
				(SELECT name from "order_status" WHERE id = status_id) as status_name,
				total_amount,
				currency,
				payment_id,
				delivery_id,
				created_at,
				updated_at`).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entity.Order{}, repository.ErrOrderAlreadyExists
			}
		}
		logrus.Errorf("OrderRepository.Create: query error: %v", err)
		return entity.Order{}, err
	}

	rowOrder, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[RowOrder])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logrus.Errorf("OrderRepository.Create: no rows returned after insert: %v", err)
			return entity.Order{}, err
		}
		logrus.Errorf("OrderRepository.Create: scan error: %v", err)
		return entity.Order{}, err
	}

	logrus.Infof("OrderRepository.CreateOrder: created orderID=%v for customerID=%v", order.ID, order.CustomerID)

	return rowOrder.ToEntity(), nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.StatusName, time time.Time) error {
	logrus.Infof("OrderRepository.UpdateOrderStatus: orderID=%v, status=%v", orderID, status)

	query, args, _ := r.Builder.
		Update("orders").
		Set("status", squirrel.Expr("(SELECT id FROM order_status WHERE name = ?)", string(status))).
		Set("updated_at", time).
		Where("id = ?", orderID).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)

	if err != nil {
		logrus.Errorf("OrderRepository.UpdateOrderStatus: failed to update order status: %v", err)
		return repository.ErrCannotUpdateOrder
	}

	logrus.Infof("OrderRepository.UpdateOrderStatus: updated orderID=%v to status=%v", orderID, status)
	return nil
}

func (r *Repository) UpdateOrderPayment(ctx context.Context, orderID, paymentID uuid.UUID, time time.Time) error {
	logrus.Infof("OrderRepository.UpdateOrderPayment: orderID=%v", orderID)

	query, args, _ := r.Builder.
		Update("orders").
		Set("payment_id", paymentID). // Set status to Paid
		Set("status", squirrel.Expr("(SELECT id FROM order_status WHERE name = ?)", string(entity.StatusPaid))).
		Set("updated_at", time).
		Where("id = ?", orderID).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)

	if err != nil {
		logrus.Errorf("OrderRepository.UpdateOrderPayment: failed to update order payment: %v", err)
		return repository.ErrCannotUpdateOrder
	}

	logrus.Infof("OrderRepository.UpdateOrderPayment: updated orderID=%v", orderID)
	return nil
}

func (r *Repository) UpdateOrderDelivery(ctx context.Context, orderID, deliveryID uuid.UUID, time time.Time) error {
	logrus.Infof("OrderRepository.UpdateOrderDelivery: orderID=%v", orderID)

	query, args, _ := r.Builder.
		Update("orders").
		Set("delivery_id", deliveryID). // Set status to Delivering
		Set("status", squirrel.Expr("(SELECT id FROM order_status WHERE name = ?)", string(entity.StatusDelivering))).
		Set("updated_at", time).
		Where("id = ?", orderID).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)

	if err != nil {
		logrus.Errorf("OrderRepository.UpdateOrderDelivery: failed to update order delivery: %v", err)
		return repository.ErrCannotUpdateOrder
	}

	logrus.Infof("OrderRepository.UpdateOrderDelivery: updated orderID=%v", orderID)
	return nil
}

func (r *Repository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (entity.Order, error) {
	logrus.Infof("OrderRepository.GetOrderByID: orderID=%v", orderID)

	query, args, _ := r.Builder.
		Select("id", "customer_id", "status", "payment_id", "delivery_id", "total_amount", "currency", "created_at", "updated_at").
		From("orders").
		Where("id = ?", orderID).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OrderRepository.GetOrderByID: query error: %v", err)
		return entity.Order{}, err
	}

	rowOrder, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[RowOrder])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logrus.Warnf("OrderRepository.GetOrderByID: order not found: %v", err)
			return entity.Order{}, repository.ErrOrderNotFound
		}
		logrus.Errorf("OrderRepository.GetOrderByID: scan error: %v", err)
		return entity.Order{}, err
	}

	order := rowOrder.ToEntity()

	query, args, _ = r.Builder.
		Select("id, product_id, product_name, product_price, amount, total_price, notes").
		From("order_items").
		Where("order_id = ?", orderID).
		ToSql()

	rows, err = r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OrderRepository.GetOrderByID: items query error: %v", err)
		return entity.Order{}, err
	}

	rowsItem, err := pgx.CollectRows(rows, pgx.RowToStructByName[RowItem])
	if err != nil {
		logrus.Errorf("OrderRepository.GetOrderByID: items scan error: %v", err)
		return entity.Order{}, err
	}

	order.Items = lo.Map(rowsItem, func(r RowItem, _ int) entity.OrderItem { return r.ToEntity() })

	logrus.Infof("OrderRepository.GetOrderByID: retrieved orderID=%v with %d items", order.ID, len(order.Items))
	return order, nil
}

// Return page of found orders sorted by creation time, total items, and error.
func (r *Repository) GetAllOrders(ctx context.Context, limit, offset int) (orders []entity.Order, total int, err error) {
	logrus.Infof("OrderRepository.GetAllOrders: limit=%d offset=%d", limit, offset)

	// Get total
	countQuery, countArgs, _ := r.Builder.
		Select("COUNT(*)").
		From(`"orders"`).
		ToSql()

	err = r.GetTxManager(ctx).QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		logrus.Errorf("OrderRepository.GetAllOrders: count query error: %v", err)
		return nil, 0, err
	}

	// Get page of orders
	query, args, _ := r.Builder.
		Select(`
			id,
			customer_id,
			status_id,
			(SELECT name from order_status WHERE id = status_id) as status_name,
			total_amount,
			currency,
			payment_id,
			delivery_id,
			created_at,
			updated_at
		`).
		From(`"orders"`).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OrderRepository.GetAllOrders: orders query error: %v", err)
		return nil, 0, err
	}

	rowOrders, err := pgx.CollectRows(rows, pgx.RowToStructByName[RowOrder])
	if err != nil {
		logrus.Errorf("OrderRepository.GetAllOrders: orders scan error: %v", err)
		return nil, 0, err
	}

	if len(rowOrders) == 0 {
		return []entity.Order{}, total, nil
	}

	var orderIDs []uuid.UUID

	// Convert orders to entity + collect IDs
	orders = lo.Map(rowOrders, func(r RowOrder, _ int) entity.Order {
		orderIDs = append(orderIDs, r.ID)
		return r.ToEntity()
	})

	// Get items for orders
	itemsQuery, itemsArgs, _ := r.Builder.
		Select(`
			id,
			order_id,
			product_id,
			product_name,
			product_price,
			amount,
			total_price,
			notes
		`).
		From("order_item").
		Where("order_id IN (?)", orderIDs).
		ToSql()

	itemRows, err := r.GetTxManager(ctx).Query(ctx, itemsQuery, itemsArgs...)
	if err != nil {
		logrus.Errorf("OrderRepository.GetAllOrders: items query error: %v", err)
		return nil, 0, err
	}

	rowItems, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[RowItem])
	if err != nil {
		logrus.Errorf("OrderRepository.GetAllOrders: items scan error: %v", err)
		return nil, 0, err
	}

	// Group items by OrderID
	itemsByOrder := make(map[uuid.UUID][]entity.OrderItem)
	for _, r := range rowItems {
		item := r.ToEntity()
		itemsByOrder[r.OrderID] = append(itemsByOrder[r.OrderID], item)
	}

	// Add items to orders
	for i := range orders {
		orders[i].Items = itemsByOrder[orders[i].ID]
	}

	logrus.Infof("OrderRepository.GetAllOrders: fetched %d orders", len(orders))
	return orders, total, nil
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) (orders []entity.Order, err error) {
	logrus.Infof("OrderRepository.GetOrdersByUserID: userID = %v", userID)

	// Get page of orders
	query, args, _ := r.Builder.
		Select(`
			id,
			customer_id,
			status_id,
			(SELECT name from order_status WHERE id = status_id) as status_name,
			total_amount,
			currency,
			payment_id,
			delivery_id,
			created_at,
			updated_at
		`).
		From(`"orders"`).
		OrderBy("created_at DESC").
		Where("customer_id = ?", userID).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OrderRepository.GetOrdersByUserID: orders query error: %v", err)
		return nil, err
	}

	rowOrders, err := pgx.CollectRows(rows, pgx.RowToStructByName[RowOrder])
	if err != nil {
		logrus.Errorf("OrderRepository.GetOrdersByUserID: orders scan error: %v", err)
		return nil, err
	}

	if len(rowOrders) == 0 {
		return []entity.Order{}, nil
	}

	var orderIDs []uuid.UUID

	// Convert orders to entity + collect IDs
	orders = lo.Map(rowOrders, func(r RowOrder, _ int) entity.Order {
		orderIDs = append(orderIDs, r.ID)
		return r.ToEntity()
	})

	// Get items for orders
	itemsQuery, itemsArgs, _ := r.Builder.
		Select(`
			id,
			order_id,
			product_id,
			product_name,
			product_price,
			amount,
			total_price,
			notes
		`).
		From("order_item").
		Where("order_id IN (?)", orderIDs).
		ToSql()

	itemRows, err := r.GetTxManager(ctx).Query(ctx, itemsQuery, itemsArgs...)
	if err != nil {
		logrus.Errorf("OrderRepository.GetOrdersByUserID: items query error: %v", err)
		return nil, err
	}

	rowItems, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[RowItem])
	if err != nil {
		logrus.Errorf("OrderRepository.GetOrdersByUserID: items scan error: %v", err)
		return nil, err
	}

	// Group items by OrderID
	itemsByOrder := make(map[uuid.UUID][]entity.OrderItem)
	for _, r := range rowItems {
		item := r.ToEntity()
		itemsByOrder[r.OrderID] = append(itemsByOrder[r.OrderID], item)
	}

	// Add items to orders
	for i := range orders {
		orders[i].Items = itemsByOrder[orders[i].ID]
	}

	logrus.Infof("OrderRepository.GetOrdersByUserID: fetched %d orders", len(orders))
	return orders, nil
}
