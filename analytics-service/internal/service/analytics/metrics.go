package analytics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics содержит Prometheus метрики для аналитики
type Metrics struct {
	OrderEventsTotal  *prometheus.CounterVec
	OrderAmount       *prometheus.HistogramVec
}

// NewMetrics создает новый экземпляр метрик
func NewMetrics() *Metrics {
	return &Metrics{
		OrderEventsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "order_events_total",
				Help: "Total number of order events",
			},
			[]string{"event_type"},
		),
		OrderAmount: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "order_amount",
				Help:    "Order amounts",
				Buckets: prometheus.ExponentialBuckets(100, 2, 10), // 100, 200, 400, 800, ..., 51200
			},
			[]string{"event_type"},
		),
	}
}

// RecordOrderEvent увеличивает счетчик событий
func (m *Metrics) RecordOrderEvent(eventType string) {
	m.OrderEventsTotal.WithLabelValues(eventType).Inc()
}

// RecordOrderAmount записывает сумму заказа
func (m *Metrics) RecordOrderAmount(amount float64, eventType string) {
	m.OrderAmount.WithLabelValues(eventType).Observe(amount)
}

