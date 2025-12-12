package app

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/subscription-service/pkg/validator"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	for _, r := range handler.Routes() {
		fmt.Printf("%s %s\n", r.Method, r.Path)
	}

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {

	orderGroup := handler.Group("orders")
	{
		orderGroup.POST("", app.PostOrderHandler().Handle)
		orderGroup.GET("", app.GetAllOrdersHandler().Handle)
		orderGroup.GET("/:id", app.GetOrderHandler().Handle)
		orderGroup.GET("/user/:userId", app.GetOrdersByUserHandler().Handle)
		orderGroup.GET("/user/:userId/active", app.GetActiveOrdersByUserHandler().Handle)
	}

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	// Swagger UI
	handler.GET("/swagger/*", echoSwagger.WrapHandler)
}
