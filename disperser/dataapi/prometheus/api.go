package prometheus

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promconfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

var (
	clientOnce sync.Once
	apiIntance *prometheusApi
)

type Api interface {
	QueryRange(ctx context.Context, query string, start time.Time, end time.Time, step time.Duration) (model.Value, v1.Warnings, error)
}

type prometheusApi struct {
	api v1.API
}

var _ Api = (*prometheusApi)(nil)

func NewApi(config Config) (*prometheusApi, error) {
	var err error
	clientOnce.Do(func() {
		roundTripper := promconfig.NewBasicAuthRoundTripper(config.Username, promconfig.Secret(config.Secret), "", api.DefaultRoundTripper)
		client, errN := api.NewClient(api.Config{
			Address:      config.ServerURL,
			RoundTripper: roundTripper,
		})
		if errN != nil {
			err = errN
			return
		}
		v1api := v1.NewAPI(client)
		apiIntance = &prometheusApi{
			api: v1api,
		}
	})

	return apiIntance, err
}

func (p *prometheusApi) QueryRange(
	ctx context.Context,
	query string,
	start time.Time,
	end time.Time,
	step time.Duration,
) (model.Value, v1.Warnings, error) {
	result, warnings, err := p.api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	})
	if err != nil {
		return nil, nil, err
	}
	return result, warnings, nil
}
