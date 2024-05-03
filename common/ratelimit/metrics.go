package ratelimit

import "github.com/prometheus/client_golang/prometheus"

func RegisterMetrics(registerer prometheus.Registerer) {
	registerer.MustRegister(prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rate_limiter_bucket_levels",
		Help: "Current level of each bucket for rate limiting",
	}, []string{"requester_id", "bucket_index"}))
}
