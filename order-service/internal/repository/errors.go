package repository

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrCannotCreateOrder  = errors.New("cannot create order")
	ErrCannotUpdateOrder  = errors.New("cannot update order")
	ErrOrderNotFound      = errors.New("order not found")
)
