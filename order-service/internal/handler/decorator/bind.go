package decorator

import (
	"errors"
	"net/http"

	h "github.com/4udiwe/big-bob-pizza/order-service/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type handler[T any] interface {
	Handle(c echo.Context, in T) error
}

type bindAndValidateDecorator[T any] struct {
	inner handler[T]
}

func NewBindAndValidateDecorator[T any](inner handler[T]) h.Handler {
	return &bindAndValidateDecorator[T]{inner: inner}
}

func (d *bindAndValidateDecorator[T]) Handle(c echo.Context) error {
	logrus.Infof("HTTP %s %s from %s", c.Request().Method, c.Path(), c.Request().RemoteAddr)

	var in T

	if err := c.Bind(&in); err != nil {
		logrus.Errorf("Failed to bind request: %v", err)
		return d.handleError(err, err.Error())
	}

	if err := c.Validate(in); err != nil {
		logrus.Errorf("Failed to validate request: %v", err)
		return d.handleError(err, err.Error())
	}

	return d.inner.Handle(c, in)
}

func (d *bindAndValidateDecorator[T]) handleError(err error, defaultMsg string) *echo.HTTPError {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return echo.NewHTTPError(http.StatusBadRequest, httpErr.Message)
	}
	return echo.NewHTTPError(http.StatusBadRequest, defaultMsg)
}
