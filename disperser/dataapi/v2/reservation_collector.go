package v2

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
)

// ReservationExpirationCollector is a custom Prometheus collector that queries reservation data
// and exposes metrics about expiring reservations
type ReservationExpirationCollector struct {
	subgraphClient dataapi.SubgraphClient
	logger         logging.Logger

	// Metrics
	reservationsActive         prometheus.Gauge
	reservationTimeUntilExpiry *prometheus.GaugeVec
}

// NewReservationExpirationCollector creates a new collector
func NewReservationExpirationCollector(subgraphClient dataapi.SubgraphClient, metrics *dataapi.Metrics, logger logging.Logger) *ReservationExpirationCollector {
	return &ReservationExpirationCollector{
		subgraphClient: subgraphClient,
		logger:         logger,
		reservationsActive: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "eigenda_reservations_active",
			Help: "Number of active reservations",
		}),
		reservationTimeUntilExpiry: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "eigenda_reservation_time_until_expiry_seconds",
			Help: "Time until reservation expiration in seconds",
		}, []string{"account"}),
	}
}

// Describe implements prometheus.Collector
func (c *ReservationExpirationCollector) Describe(ch chan<- *prometheus.Desc) {
	c.reservationsActive.Describe(ch)
	c.reservationTimeUntilExpiry.Describe(ch)
}

// Collect implements prometheus.Collector
func (c *ReservationExpirationCollector) Collect(ch chan<- prometheus.Metric) {
	// Update counts with timeout to prevent blocking Prometheus scrapes
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	c.updateCounts(ctx)

	// Collect metrics
	c.reservationsActive.Collect(ch)
	c.reservationTimeUntilExpiry.Collect(ch)
}

// updateCounts queries the GraphQL endpoint and updates the metrics
func (c *ReservationExpirationCollector) updateCounts(ctx context.Context) {
	// Query all active reservations
	currentTimestamp := uint64(time.Now().Unix())
	reservations, err := c.subgraphClient.QueryReservations(ctx, currentTimestamp, 1000, 0)
	if err != nil {
		c.logger.Warn("Failed to query reservations", "error", err)
		return
	}

	// Calculate metrics
	now := time.Now()
	activeCount := 0
	expiringCounts := map[string]int{
		"24h": 0,
		"7d":  0,
		"3m":  0,
	}

	// Clear metrics before adding new observations
	c.reservationTimeUntilExpiry.Reset()

	for _, res := range reservations {
		// Calculate time until expiration
		expirationTime := time.Unix(int64(res.EndTimestamp), 0)
		timeUntilExpiration := expirationTime.Sub(now)

		// Skip already expired reservations
		if timeUntilExpiration < 0 {
			continue
		}

		activeCount++

		// Count expiring reservations by time window
		if timeUntilExpiration <= 24*time.Hour {
			expiringCounts["24h"]++
		} else if timeUntilExpiration <= 7*24*time.Hour {
			expiringCounts["7d"]++
		} else if timeUntilExpiration <= 3*30*24*time.Hour {
			expiringCounts["3m"]++
		} else {
			// Ignore reservations that are not expiring in the next 3 months
			continue
		}

		// Record gauge value
		c.reservationTimeUntilExpiry.WithLabelValues(string(res.Account)).Set(timeUntilExpiration.Seconds())
	}

	// Update gauges
	c.reservationsActive.Set(float64(activeCount))

	c.logger.Info("Updated reservation metrics", "active", activeCount, "expiring_24h", expiringCounts["24h"], "expiring_7d", expiringCounts["7d"], "expiring_3m", expiringCounts["3m"])
}
