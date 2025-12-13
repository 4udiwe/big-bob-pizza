package order_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	mock_transactor "github.com/4udiwe/big-bob-pizza/order-service/internal/mocks"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/repository"
	service "github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/service/order/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestService_CreateOrder(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	orderID := uuid.New()

	order := entity.Order{
		CustomerID:  customerID,
		TotalAmount: 100.0,
		Currency:    "USD",
		Items: []entity.OrderItem{
			{
				ProductID:    uuid.New(),
				ProductName:  "Pizza",
				ProductPrice: 50.0,
				Amount:       2,
				TotalPrice:   100.0,
			},
		},
	}

	tests := []struct {
		name        string
		setup       func(orderRepo *mocks.MockOrderRepo, itemsRepo *mocks.MockItemsRepo, outboxRepo *mocks.MockOutboxRepo, cacheRepo *mocks.MockCacheRepo, tx *mock_transactor.MockTransactor)
		expectedErr error
	}{
		{
			name: "order already exists",
			setup: func(orderRepo *mocks.MockOrderRepo, itemsRepo *mocks.MockItemsRepo, outboxRepo *mocks.MockOutboxRepo, cacheRepo *mocks.MockCacheRepo, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().
					WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				orderRepo.EXPECT().
					Create(gomock.Any(), order).
					Return(entity.Order{}, repository.ErrOrderAlreadyExists)
			},
			expectedErr: service.ErrOrderAlreadyExists,
		},
		{
			name: "create order error",
			setup: func(orderRepo *mocks.MockOrderRepo, itemsRepo *mocks.MockItemsRepo, outboxRepo *mocks.MockOutboxRepo, cacheRepo *mocks.MockCacheRepo, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().
					WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				orderRepo.EXPECT().
					Create(gomock.Any(), order).
					Return(entity.Order{}, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
		{
			name: "insert items error",
			setup: func(orderRepo *mocks.MockOrderRepo, itemsRepo *mocks.MockItemsRepo, outboxRepo *mocks.MockOutboxRepo, cacheRepo *mocks.MockCacheRepo, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().
					WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				createdOrder := entity.Order{
					ID:          orderID,
					CustomerID:  customerID,
					TotalAmount: 100.0,
					Currency:    "USD",
					Status:      entity.OrderStatus{Name: entity.StatusCreated},
				}

				orderRepo.EXPECT().
					Create(gomock.Any(), order).
					Return(createdOrder, nil)

				itemsRepo.EXPECT().
					InsertItems(gomock.Any(), orderID, order.Items).
					Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
		{
			name: "success",
			setup: func(orderRepo *mocks.MockOrderRepo, itemsRepo *mocks.MockItemsRepo, outboxRepo *mocks.MockOutboxRepo, cacheRepo *mocks.MockCacheRepo, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().
					WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				createdOrder := entity.Order{
					ID:          orderID,
					CustomerID:  customerID,
					TotalAmount: 100.0,
					Currency:    "USD",
					Status:      entity.OrderStatus{Name: entity.StatusCreated},
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				items := []entity.OrderItem{
					{
						ID:           uuid.New(),
						ProductID:    order.Items[0].ProductID,
						ProductName:  order.Items[0].ProductName,
						ProductPrice: order.Items[0].ProductPrice,
						Amount:       order.Items[0].Amount,
						TotalPrice:   order.Items[0].TotalPrice,
					},
				}

				orderRepo.EXPECT().
					Create(gomock.Any(), order).
					Return(createdOrder, nil)

				itemsRepo.EXPECT().
					InsertItems(gomock.Any(), orderID, order.Items).
					Return(items, nil)

				outboxRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				createdOrder.Items = items
				cacheRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(nil)

				cacheRepo.EXPECT().
					AddToActive(gomock.Any(), gomock.Any()).
					Return(nil)

				cacheRepo.EXPECT().
					AddUserActive(gomock.Any(), customerID, orderID).
					Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := mocks.NewMockOrderRepo(ctrl)
			itemsRepo := mocks.NewMockItemsRepo(ctrl)
			outboxRepo := mocks.NewMockOutboxRepo(ctrl)
			cacheRepo := mocks.NewMockCacheRepo(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			svc := service.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, tx)

			tt.setup(orderRepo, itemsRepo, outboxRepo, cacheRepo, tx)

			_, err := svc.CreateOrder(ctx, order)
			if !errors.Is(err, tt.expectedErr) && (tt.expectedErr == nil || err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestService_GetOrderByID(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	customerID := uuid.New()

	order := entity.Order{
		ID:          orderID,
		CustomerID:  customerID,
		TotalAmount: 100.0,
		Currency:    "USD",
		Status:      entity.OrderStatus{Name: entity.StatusCreated},
		Items: []entity.OrderItem{
			{
				ID:           uuid.New(),
				ProductID:    uuid.New(),
				ProductName:  "Pizza",
				ProductPrice: 50.0,
				Amount:       2,
				TotalPrice:   100.0,
			},
		},
	}

	tests := []struct {
		name        string
		setup       func(cacheRepo *mocks.MockCacheRepo, orderRepo *mocks.MockOrderRepo)
		expectedErr error
	}{
		{
			name: "cache hit",
			setup: func(cacheRepo *mocks.MockCacheRepo, orderRepo *mocks.MockOrderRepo) {
				cacheRepo.EXPECT().
					GetByID(gomock.Any(), orderID).
					Return(&order, nil)
			},
			expectedErr: nil,
		},
		{
			name: "cache miss, fallback to postgres",
			setup: func(cacheRepo *mocks.MockCacheRepo, orderRepo *mocks.MockOrderRepo) {
				cacheRepo.EXPECT().
					GetByID(gomock.Any(), orderID).
					Return(nil, errors.New("not found"))

				orderRepo.EXPECT().
					GetOrderByID(gomock.Any(), orderID).
					Return(order, nil)

				cacheRepo.EXPECT().
					Save(gomock.Any(), &order).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "order not found",
			setup: func(cacheRepo *mocks.MockCacheRepo, orderRepo *mocks.MockOrderRepo) {
				cacheRepo.EXPECT().
					GetByID(gomock.Any(), orderID).
					Return(nil, errors.New("not found"))

				orderRepo.EXPECT().
					GetOrderByID(gomock.Any(), orderID).
					Return(entity.Order{}, repository.ErrOrderNotFound)
			},
			expectedErr: repository.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cacheRepo := mocks.NewMockCacheRepo(ctrl)
			orderRepo := mocks.NewMockOrderRepo(ctrl)
			itemsRepo := mocks.NewMockItemsRepo(ctrl)
			outboxRepo := mocks.NewMockOutboxRepo(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			svc := service.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, tx)

			tt.setup(cacheRepo, orderRepo)

			_, err := svc.GetOrderByID(ctx, orderID)
			if !errors.Is(err, tt.expectedErr) && (tt.expectedErr == nil || err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestService_GetOrdersByUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name        string
		setup       func(orderRepo *mocks.MockOrderRepo)
		expectedErr error
	}{
		{
			name: "repo error",
			setup: func(orderRepo *mocks.MockOrderRepo) {
				orderRepo.EXPECT().
					GetOrdersByUserID(gomock.Any(), userID, 20, 0).
					Return(nil, 0, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
		{
			name: "success",
			setup: func(orderRepo *mocks.MockOrderRepo) {
				orders := []entity.Order{
					{
						ID:          uuid.New(),
						CustomerID:  userID,
						TotalAmount: 100.0,
						Currency:    "USD",
						Status:      entity.OrderStatus{Name: entity.StatusCreated},
					},
				}
				orderRepo.EXPECT().
					GetOrdersByUserID(gomock.Any(), userID, 20, 0).
					Return(orders, 1, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := mocks.NewMockOrderRepo(ctrl)
			itemsRepo := mocks.NewMockItemsRepo(ctrl)
			outboxRepo := mocks.NewMockOutboxRepo(ctrl)
			cacheRepo := mocks.NewMockCacheRepo(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			svc := service.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, tx)

			tt.setup(orderRepo)

			_, _, err := svc.GetOrdersByUser(ctx, userID, 20, 0)
			if !errors.Is(err, tt.expectedErr) && (tt.expectedErr == nil || err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestService_GetAllOrders(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setup       func(orderRepo *mocks.MockOrderRepo)
		expectedErr error
	}{
		{
			name: "repo error",
			setup: func(orderRepo *mocks.MockOrderRepo) {
				orderRepo.EXPECT().
					GetAllOrders(gomock.Any(), 20, 0).
					Return(nil, 0, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
		{
			name: "success",
			setup: func(orderRepo *mocks.MockOrderRepo) {
				orders := []entity.Order{
					{
						ID:          uuid.New(),
						CustomerID:  uuid.New(),
						TotalAmount: 100.0,
						Currency:    "USD",
						Status:      entity.OrderStatus{Name: entity.StatusCreated},
					},
				}
				orderRepo.EXPECT().
					GetAllOrders(gomock.Any(), 20, 0).
					Return(orders, 1, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := mocks.NewMockOrderRepo(ctrl)
			itemsRepo := mocks.NewMockItemsRepo(ctrl)
			outboxRepo := mocks.NewMockOutboxRepo(ctrl)
			cacheRepo := mocks.NewMockCacheRepo(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			svc := service.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, tx)

			tt.setup(orderRepo)

			_, _, err := svc.GetAllOrders(ctx, 20, 0)
			if !errors.Is(err, tt.expectedErr) && (tt.expectedErr == nil || err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestService_GetActiveOrdersByUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()

	order := entity.Order{
		ID:          orderID,
		CustomerID:  userID,
		TotalAmount: 100.0,
		Currency:    "USD",
		Status:      entity.OrderStatus{Name: entity.StatusCreated},
	}

	tests := []struct {
		name        string
		setup       func(cacheRepo *mocks.MockCacheRepo)
		expectedErr error
	}{
		{
			name: "no active orders",
			setup: func(cacheRepo *mocks.MockCacheRepo) {
				cacheRepo.EXPECT().
					GetUserActiveOrders(gomock.Any(), userID).
					Return([]string{}, nil)
			},
			expectedErr: service.ErrNoActiveOrders,
		},
		{
			name: "success with active orders",
			setup: func(cacheRepo *mocks.MockCacheRepo) {
				cacheRepo.EXPECT().
					GetUserActiveOrders(gomock.Any(), userID).
					Return([]string{orderID.String()}, nil)

				cacheRepo.EXPECT().
					GetByID(gomock.Any(), orderID).
					Return(&order, nil)
			},
			expectedErr: nil,
		},
		{
			name: "cache error",
			setup: func(cacheRepo *mocks.MockCacheRepo) {
				cacheRepo.EXPECT().
					GetUserActiveOrders(gomock.Any(), userID).
					Return(nil, errors.New("cache error"))
			},
			expectedErr: service.ErrNoActiveOrders,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := mocks.NewMockOrderRepo(ctrl)
			itemsRepo := mocks.NewMockItemsRepo(ctrl)
			outboxRepo := mocks.NewMockOutboxRepo(ctrl)
			cacheRepo := mocks.NewMockCacheRepo(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			svc := service.NewService(orderRepo, itemsRepo, outboxRepo, cacheRepo, tx)

			tt.setup(cacheRepo)

			_, err := svc.GetActiveOrdersByUser(ctx, userID)
			if !errors.Is(err, tt.expectedErr) && (tt.expectedErr == nil || err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}


