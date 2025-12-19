package app

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/subscription-service/pkg/validator"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	analyticsGroup := handler.Group("analytics")
	{
		analyticsGroup.GET("/stats", app.GetStatsHandler().Handle)
		analyticsGroup.GET("/revenue", app.GetRevenueHandler().Handle)
		analyticsGroup.GET("/orders/:orderId/events", app.GetOrderEventsHandler().Handle)
	}

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	// Prometheus metrics endpoint
	if app.cfg.Prometheus.Enabled {
		handler.GET(app.cfg.Prometheus.Path, func(c echo.Context) error {
			promhttp.Handler().ServeHTTP(c.Response(), c.Request())
			return nil
		})
	}

	// Swagger UI
	handler.GET("/swagger/*", echoSwagger.WrapHandler)
}

