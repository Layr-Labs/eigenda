package main

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/metrics"
	"math/rand"
	"time"
)

// TODO don't merge this, this is just a test bed

func main() {

	metricsConfig := &metrics.Config{
		Namespace: "test",
		HTTPPort:  9101,
	}

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		panic(err)
	}

	metricsServer := metrics.NewMetrics(logger, metricsConfig)

	l1, err := metricsServer.NewLatencyMetric(
		"l1",
		"",
		metrics.NewQuantile(0.5),
		metrics.NewQuantile(0.9),
		metrics.NewQuantile(0.99))
	if err != nil {
		panic(err)
	}

	l1HALF, err := metricsServer.NewLatencyMetric(
		"l1",
		"HALF",
		metrics.NewQuantile(0.5),
		metrics.NewQuantile(0.9),
		metrics.NewQuantile(0.99))
	if err != nil {
		panic(err)
	}

	c1, err := metricsServer.NewCountMetric("c1", "")
	if err != nil {
		panic(err)
	}

	c1DOUBLE, err := metricsServer.NewCountMetric("c1", "DOUBLE")
	if err != nil {
		panic(err)
	}

	metricsServer.Start()

	prev := time.Now()
	for i := 0; i < 100000; i++ {
		fmt.Printf("Iteration %d\n", i)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		now := time.Now()
		elapsed := now.Sub(prev)
		prev = now

		l1.ReportLatency(elapsed)
		l1HALF.ReportLatency(elapsed / 2)

		c1.Increment()
		c1DOUBLE.Add(2)
	}

}
