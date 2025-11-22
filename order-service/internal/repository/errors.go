package repository

import "errors"

var (
	ErrCannotCreateOrder = errors.New("cannot create order")
	ErrCannotUpdateOrder = errors.New("cannot update order")
	
)
