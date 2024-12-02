package main

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/metrics"
	"math/rand"
	"sync/atomic"
	"time"
)

// This is a simple test bed for validating the metrics server (since it's not straight forward to unit test).

type LabelType1 struct {
	foo string
	bar string
	baz string
}

type LabelType2 struct {
	X string
	Y string
	Z string
}

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
		"this metric shows the latency of the sleep cycle",
		LabelType1{},
		metrics.NewQuantile(0.5),
		metrics.NewQuantile(0.9),
		metrics.NewQuantile(0.99))
	if err != nil {
		panic(err)
	}

	c1, err := metricsServer.NewCountMetric(
		"c1",
		"this metric shows the number of times the sleep cycle has been executed",
		LabelType2{})
	if err != nil {
		panic(err)
	}

	c2, err := metricsServer.NewCountMetric(
		"c2",
		"the purpose of this counter is to test what happens if we don't provide a label template",
		nil)
	if err != nil {
		panic(err)
	}

	g1, err := metricsServer.NewGaugeMetric(
		"g1",
		"milliseconds",
		"this metric shows the duration of the most recent sleep cycle",
		LabelType1{})
	if err != nil {
		panic(err)
	}

	sum := atomic.Int64{}
	err = metricsServer.NewAutoGauge(
		"g2",
		"milliseconds",
		"this metric shows the sum of all sleep cycles",
		1*time.Second,
		func() float64 {
			return float64(sum.Load())
		},
		LabelType2{X: "sum"})
	if err != nil {
		panic(err)
	}

	err = metricsServer.WriteMetricsDocumentation("metrics.md")
	if err != nil {
		panic(err)
	}

	err = metricsServer.Start()
	if err != nil {
		panic(err)
	}

	prev := time.Now()
	for i := 0; i < 100; i++ {
		fmt.Printf("Iteration %d\n", i)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		now := time.Now()
		elapsed := now.Sub(prev)
		prev = now

		l1.ReportLatency(elapsed)

		l1.ReportLatency(elapsed/2,
			LabelType1{
				foo: "half of the normal value",
				bar: "42",
				baz: "true",
			})

		c1.Increment()
		c1.Add(2, LabelType2{
			X: "2x",
		})
		c2.Increment()

		g1.Set(float64(elapsed.Milliseconds()),
			LabelType1{
				foo: "bar",
				bar: "baz",
				baz: "foo",
			})

		sum.Store(sum.Load() + elapsed.Milliseconds())
	}

	err = metricsServer.Stop()
	if err != nil {
		panic(err)
	}
}
