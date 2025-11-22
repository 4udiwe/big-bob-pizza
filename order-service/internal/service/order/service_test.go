package order_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	mock_transactor "github.com/4udiwe/big-bob-pizza/order-service/internal/mocks"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/service/order/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()

	customerID := uuid.New()
	amount := 777.7
	currency := "RUB"

	createdOrder := entity.Product{
		ID:          uuid.New(),
		CustomerID:  customerID,
		Status:      entity.OrderStatusCreated,
		TotalAmount: amount,
		Currency:    currency,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	arbitraryErr := errors.New("arbitrary error")

	type mockBehavior func(r *mocks.MockOrderRepository)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		want         entity.Product
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					CreateOrder(ctx, customerID, amount, currency).
					Return(createdOrder, nil).
					Times(1)
			},
			want: createdOrder,
		},
		{
			name: "repo error",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					CreateOrder(ctx, customerID, amount, currency).
					Return(entity.Product{}, arbitraryErr).
					Times(1)
			},
			want:    entity.Product{},
			wantErr: arbitraryErr,
		},
		{
			name: "empty order returned",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					CreateOrder(ctx, customerID, amount, currency).
					Return(entity.Product{}, nil).
					Times(1)
			},
			want:    entity.Product{},
			wantErr: order.ErrCannotCreateOrder,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockOrderRepository(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(repo)

			s := order.New(repo, tx)

			got, err := s.CreateOrder(ctx, customerID, amount, currency)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	ctx := context.Background()

	orderID := uuid.New()
	status := entity.OrderStatusPaid
	arbitraryErr := errors.New("arbitrary error")

	type mockBehavior func(repo *mocks.MockOrderRepository, tx *mock_transactor.MockTransactor)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(repo *mocks.MockOrderRepository, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)

				repo.EXPECT().
					UpdateOrderStatus(ctx, orderID, status).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name: "repo error",
			mockBehavior: func(repo *mocks.MockOrderRepository, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)

				repo.EXPECT().
					UpdateOrderStatus(ctx, orderID, status).
					Return(arbitraryErr).
					Times(1)
			},
			wantErr: order.ErrCannotUpdateOrder,
		},
		{
			name: "transaction fails",
			mockBehavior: func(repo *mocks.MockOrderRepository, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return arbitraryErr
					}).
					Times(1)

				repo.EXPECT().
					UpdateOrderStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: order.ErrCannotUpdateOrder,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockOrderRepository(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(repo, tx)

			s := order.New(repo, tx)

			err := s.UpdateOrderStatus(ctx, orderID, status)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	ctx := context.Background()

	orderID := uuid.New()
	arbitraryErr := errors.New("arbitrary error")

	foundOrder := entity.Product{
		ID:          orderID,
		CustomerID:  uuid.New(),
		Status:      entity.OrderStatusPreparing,
		Currency:    "USD",
		TotalAmount: 150,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	type mockBehavior func(r *mocks.MockOrderRepository)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		want         entity.Product
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					GetOrderByID(ctx, orderID).
					Return(foundOrder, nil).
					Times(1)
			},
			want: foundOrder,
		},
		{
			name: "repo error",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					GetOrderByID(ctx, orderID).
					Return(entity.Product{}, arbitraryErr).
					Times(1)
			},
			wantErr: arbitraryErr,
		},
		{
			name: "order not found",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					GetOrderByID(ctx, orderID).
					Return(entity.Product{}, nil).
					Times(1)
			},
			wantErr: order.ErrOrderNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockOrderRepository(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(repo)

			s := order.New(repo, tx)

			got, err := s.GetOrderByID(ctx, orderID)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetAllOrders(t *testing.T) {
	ctx := context.Background()

	orders := []entity.Product{
		{
			ID:          uuid.New(),
			CustomerID:  uuid.New(),
			Status:      entity.OrderStatusPaid,
			TotalAmount: 300,
			Currency:    "RUB",
		},
		{
			ID:          uuid.New(),
			CustomerID:  uuid.New(),
			Status:      entity.OrderStatusDelivering,
			TotalAmount: 500,
			Currency:    "USD",
		},
	}

	arbitraryErr := errors.New("arbitrary error")

	type mockBehavior func(r *mocks.MockOrderRepository)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		want         []entity.Product
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					GetAllOrders(ctx).
					Return(orders, nil).
					Times(1)
			},
			want: orders,
		},
		{
			name: "repo error",
			mockBehavior: func(r *mocks.MockOrderRepository) {
				r.EXPECT().
					GetAllOrders(ctx).
					Return(nil, arbitraryErr).
					Times(1)
			},
			want:    nil,
			wantErr: order.ErrCannotFetchOrders,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockOrderRepository(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(repo)

			s := order.New(repo, tx)

			got, err := s.GetAllOrders(ctx)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
