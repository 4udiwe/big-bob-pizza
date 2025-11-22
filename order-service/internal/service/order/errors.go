package order

import "errors"

var (
	ErrCannotCreateOrder = errors.New("cannot create order")
	ErrCannotUpdateOrder = errors.New("cannot update order")
	ErrOrderNotFound     = errors.New("order not found")
	ErrCannotFetchOrders = errors.New("cannot fetch orders")
)
