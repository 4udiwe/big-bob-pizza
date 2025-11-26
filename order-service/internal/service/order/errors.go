package order

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrNoActiveOrders     = errors.New("user has no active orders")
)
