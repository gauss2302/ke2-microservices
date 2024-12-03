package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	ProductsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "product_created_total",
			Help: "Total number of products created",
		},
	)

	ProductsDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "product_deleted_total",
			Help: "Total number of products deleted",
		},
	)

	ProductsUpdatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "product_updated_total",
			Help: "Total number of products updated",
		},
	)
)
