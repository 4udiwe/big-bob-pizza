package app

import (
	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/service/analytics"
)

func (app *App) AnalyticsService() *analytics.Service {
	return app.analyticsService
}

