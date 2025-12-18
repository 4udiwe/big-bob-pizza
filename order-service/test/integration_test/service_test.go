//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	cache_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/cache"
	item_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/item"
	order_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/order"
	outbox_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/outbox"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CreateOrder_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup repositories
	orderRepo := order_repository.New(testPostgres)
	itemsRepo := item_repository.New(testPostgres)
	outboxRepo := outbox_repository.New(testPostgres)
	cacheRepo := cache_repository.NewCacheOrderRepository(testRedis)
	txManager := testPostgres

	svc := order.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, txManager)

	customerID := uuid.New()
	order := entity.Order{
		CustomerID:  customerID,
		TotalAmount: 150.0,
		Currency:    "USD",
		Items: []entity.OrderItem{
			{
				ProductID:    uuid.New(),
				ProductName:  "Pizza Margherita",
				ProductPrice: 50.0,
				Amount:       2,
				TotalPrice:   100.0,
			},
			{
				ProductID:    uuid.New(),
				ProductName:  "Coca Cola",
				ProductPrice: 5.0,
				Amount:       10,
				TotalPrice:   50.0,
			},
		},
	}

	// Create order
	created, err := svc.CreateOrder(ctx, order)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, customerID, created.CustomerID)
	assert.Equal(t, 150.0, created.TotalAmount)
	assert.Equal(t, "USD", created.Currency)
	assert.Equal(t, entity.StatusCreated, created.Status.Name)
	assert.Len(t, created.Items, 2)

	// Verify items
	for i, item := range created.Items {
		assert.NotEqual(t, uuid.Nil, item.ID)
		assert.Equal(t, order.Items[i].ProductID, item.ProductID)
		assert.Equal(t, order.Items[i].ProductName, item.ProductName)
		assert.Equal(t, order.Items[i].Amount, item.Amount)
	}

	// Verify order can be retrieved
	retrieved, err := svc.GetOrderByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, customerID, retrieved.CustomerID)
	assert.Len(t, retrieved.Items, 2)
}

func TestService_GetOrdersByUser_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup repositories
	orderRepo := order_repository.New(testPostgres)
	itemsRepo := item_repository.New(testPostgres)
	outboxRepo := outbox_repository.New(testPostgres)
	cacheRepo := cache_repository.NewCacheOrderRepository(testRedis)
	txManager := testPostgres

	svc := order.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, txManager)

	customerID := uuid.New()

	// Create multiple orders
	for i := 0; i < 3; i++ {
		order := entity.Order{
			CustomerID:  customerID,
			TotalAmount: 100.0 + float64(i),
			Currency:    "USD",
			Items: []entity.OrderItem{
				{
					ProductID:    uuid.New(),
					ProductName:  "Product",
					ProductPrice: 50.0,
					Amount:       2,
					TotalPrice:   100.0,
				},
			},
		}
		_, err := svc.CreateOrder(ctx, order)
		require.NoError(t, err)
	}

	// Get orders by user
	orders, total, err := svc.GetOrdersByUser(ctx, customerID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 3)
	assert.GreaterOrEqual(t, len(orders), 3)

	// Verify all orders belong to the user
	for _, ord := range orders {
		assert.Equal(t, customerID, ord.CustomerID)
	}
}

func TestService_GetAllOrders_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup repositories
	orderRepo := order_repository.New(testPostgres)
	itemsRepo := item_repository.New(testPostgres)
	outboxRepo := outbox_repository.New(testPostgres)
	cacheRepo := cache_repository.NewCacheOrderRepository(testRedis)
	txManager := testPostgres

	svc := order.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, txManager)

	// Create an order
	order := entity.Order{
		CustomerID:  uuid.New(),
		TotalAmount: 200.0,
		Currency:    "USD",
		Items: []entity.OrderItem{
			{
				ProductID:    uuid.New(),
				ProductName:  "Product",
				ProductPrice: 100.0,
				Amount:       2,
				TotalPrice:   200.0,
			},
		},
	}
	_, err := svc.CreateOrder(ctx, order)
	require.NoError(t, err)

	// Get all orders
	orders, total, err := svc.GetAllOrders(ctx, 10, 0)
	require.NoError(t, err)
	assert.Greater(t, total, 0)
	assert.Greater(t, len(orders), 0)
}

func TestService_GetActiveOrdersByUser_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup repositories
	orderRepo := order_repository.New(testPostgres)
	itemsRepo := item_repository.New(testPostgres)
	outboxRepo := outbox_repository.New(testPostgres)
	cacheRepo := cache_repository.NewCacheOrderRepository(testRedis)
	txManager := testPostgres

	svc := order.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, txManager)

	customerID := uuid.New()

	// Create an active order
	order := entity.Order{
		CustomerID:  customerID,
		TotalAmount: 100.0,
		Currency:    "USD",
		Items: []entity.OrderItem{
			{
				ProductID:    uuid.New(),
				ProductName:  "Product",
				ProductPrice: 50.0,
				Amount:       2,
				TotalPrice:   100.0,
			},
		},
	}
	created, err := svc.CreateOrder(ctx, order)
	require.NoError(t, err)

	// Wait a bit for cache to be updated
	time.Sleep(100 * time.Millisecond)

	// Get active orders
	activeOrders, err := svc.GetActiveOrdersByUser(ctx, customerID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(activeOrders), 1)

	// Verify the order is in active orders
	found := false
	for _, ord := range activeOrders {
		if ord.ID == created.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "created order should be in active orders")
}

func TestService_UpdateOrderStatus_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup repositories
	orderRepo := order_repository.New(testPostgres)
	itemsRepo := item_repository.New(testPostgres)
	outboxRepo := outbox_repository.New(testPostgres)
	cacheRepo := cache_repository.NewCacheOrderRepository(testRedis)
	txManager := testPostgres

	svc := order.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, txManager)

	// Create an order
	order := entity.Order{
		CustomerID:  uuid.New(),
		TotalAmount: 100.0,
		Currency:    "USD",
		Items: []entity.OrderItem{
			{
				ProductID:    uuid.New(),
				ProductName:  "Product",
				ProductPrice: 50.0,
				Amount:       2,
				TotalPrice:   100.0,
			},
		},
	}
	created, err := svc.CreateOrder(ctx, order)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusCreated, created.Status.Name)

	// Update status to paid
	paidStatus := entity.OrderStatus{
		ID:   2,
		Name: entity.StatusPaid,
	}
	updated, err := svc.UpdateOrderStatus(ctx, created.ID, paidStatus)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusPaid, updated.Status.Name)

	// Verify in database
	retrieved, err := svc.GetOrderByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusPaid, retrieved.Status.Name)
}

