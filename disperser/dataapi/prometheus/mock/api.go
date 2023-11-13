package mock

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/mock"
)

type MockPrometheusApi struct {
	mock.Mock
}

var _ prometheus.Api = (*MockPrometheusApi)(nil)

func (m *MockPrometheusApi) QueryRange(ctx context.Context, query string, start time.Time, end time.Time, step time.Duration) (model.Value, v1.Warnings, error) {
	args := m.Called()
	var value model.Value
	if args.Get(0) != nil {
		value = args.Get(0).(model.Value)
	}
	var warnings v1.Warnings
	if args.Get(1) != nil {
		warnings = args.Get(1).(v1.Warnings)
	}
	return value, warnings, args.Error(2)
}
