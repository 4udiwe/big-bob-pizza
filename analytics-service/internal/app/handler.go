package app

import (
	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/handler"
	get_order_events "github.com/4udiwe/big-bob-pizza/analytics-service/internal/handler/get_order_events"
	get_revenue "github.com/4udiwe/big-bob-pizza/analytics-service/internal/handler/get_revenue"
	get_stats "github.com/4udiwe/big-bob-pizza/analytics-service/internal/handler/get_stats"
)

func (app *App) GetStatsHandler() handler.Handler {
	return get_stats.New(app.AnalyticsService())
}

func (app *App) GetRevenueHandler() handler.Handler {
	return get_revenue.New(app.AnalyticsService())
}

func (app *App) GetOrderEventsHandler() handler.Handler {
	return get_order_events.New(app.AnalyticsService())
}

