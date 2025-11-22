package order_repository

import (
	"context"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{Postgres: postgres}
}

func (r *Repository) CreateOrder(ctx context.Context, customerID uuid.UUID, totalAmount float64, currency string) (entity.Product, error) {
	logrus.Infof("OrderRepository.CreateOrder: customerID=%v", customerID)

	query, args, _ := r.Builder.
		Insert("order").
		Columns("customer_id", "total_amount", "currency").
		Values(customerID, totalAmount, currency).
		Suffix("RETURNING id, status, payment_id, delivery_id, created_at, updated_at").
		ToSql()

	order := entity.Product{
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Currency:    currency,
	}

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&order.ID,
		&order.Status,
		&order.PaymentID,
		&order.DeliveryID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		logrus.Errorf("OrderRepository.CreateOrder: failed to create order: %v", err)
		return entity.Product{}, repository.ErrCannotCreateOrder
	}

	logrus.Infof("OrderRepository.CreateOrder: created orderID=%v for customerID=%v", order.ID, customerID)

	return order, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) error {
	logrus.Infof("OrderRepository.UpdateOrderStatus: orderID=%v, status=%v", orderID, status)

	query, args, _ := r.Builder.
		Update("order").
		Set("status", status).
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

func (r *Repository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (entity.Product, error) {
	logrus.Infof("OrderRepository.GetOrderByID: orderID=%v", orderID)

	query, args, _ := r.Builder.
		Select("id", "customer_id", "status", "payment_id", "delivery_id", "total_amount", "currency", "created_at", "updated_at").
		From("order").
		Where("id = ?", orderID).
		ToSql()

	order := entity.Product{}

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Status,
		&order.PaymentID,
		&order.DeliveryID,
		&order.TotalAmount,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		logrus.Errorf("OrderRepository.GetOrderByID: failed to get order: %v", err)
		return entity.Product{}, err
	}

	logrus.Infof("OrderRepository.GetOrderByID: retrieved orderID=%v", order.ID)
	return order, nil
}

func (r *Repository) GetAllOrders(ctx context.Context) ([]entity.Product, error) {
	logrus.Infof("OrderRepository.GetAllOrders")
	
	query, args, _ := r.Builder.
		Select("id", "customer_id", "status", "payment_id", "delivery_id", "total_amount", "currency", "created_at", "updated_at").
		From("order").
		ToSql()
	var orders []entity.Product

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OrderRepository.GetAllOrders: failed to get orders: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order entity.Product
		err := rows.Scan(
			&order.ID,	
			&order.CustomerID,
			&order.Status,
			&order.PaymentID,
			&order.DeliveryID,
			&order.TotalAmount,
			&order.Currency,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			logrus.Errorf("OrderRepository.GetAllOrders: failed to scan order: %v", err)
			return nil, err
		}	
		orders = append(orders, order)
	}

	logrus.Infof("OrderRepository.GetAllOrders: retrieved %d orders", len(orders))
	return orders, nil

}
