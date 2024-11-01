package tracing

import (
	"context"
	"errors"
	"sync"

	ometric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"golang.org/x/sync/singleflight"
)

type OrionMetrics struct {
	guageSg singleflight.Group
	guages  sync.Map

	countersSg singleflight.Group
	counters   sync.Map

	metricProvider *metric.MeterProvider
	metricMeter    ometric.Meter

	histogramSg singleflight.Group
	histograms  sync.Map
}

func NewOrionMetrics(name string, metricProvider *metric.MeterProvider) *OrionMetrics {
	return &OrionMetrics{
		metricProvider: metricProvider,
		metricMeter:    metricProvider.Meter(name),
	}
}

func (o *OrionMetrics) Histogram(ctx context.Context, histogramName string, value float64, desc string, buckets ...float64) error {
	var err error
	ct, ok := o.histograms.Load(histogramName)
	if !ok {
		ct, err, _ = o.histogramSg.Do(histogramName, func() (interface{}, error) {
			opts := []ometric.Float64HistogramOption{}
			if len(desc) > 0 {
				opts = append(opts, ometric.WithDescription(desc))
			}
			opts = append(opts, ometric.WithExplicitBucketBoundaries(buckets...))
			c, err := o.metricMeter.Float64Histogram(histogramName, opts...)
			if err != nil {
				return nil, err
			}
			o.histograms.Store(histogramName, c)
			return c, nil
		})
	}
	if err != nil {
		return err
	}
	if ct == nil {
		return errors.New("nil guage inited")
	}

	ct.(ometric.Float64Histogram).Record(ctx, value)
	return nil
}

func (o *OrionMetrics) Guage(ctx context.Context, guageName string, value float64, desc ...string) error {
	var err error
	ct, ok := o.guages.Load(guageName)
	if !ok {
		ct, err, _ = o.guageSg.Do(guageName, func() (interface{}, error) {
			opts := []ometric.Float64GaugeOption{}
			if len(desc) > 0 {
				opts = append(opts, ometric.WithDescription(desc[0]))
			}
			c, err := o.metricMeter.Float64Gauge(guageName, opts...)
			if err != nil {
				return nil, err
			}
			o.guages.Store(guageName, c)
			return c, nil
		})
	}
	if err != nil {
		return err
	}
	if ct == nil {
		return errors.New("nil guage inited")
	}

	ct.(ometric.Float64Gauge).Record(ctx, value)
	return nil
}

func (o *OrionMetrics) Counter(ctx context.Context, counterName string, delta int64, desc ...string) error {
	var err error
	ct, ok := o.counters.Load(counterName)
	if !ok {
		ct, err, _ = o.countersSg.Do(counterName, func() (interface{}, error) {
			opts := []ometric.Int64CounterOption{}
			if len(desc) > 0 {
				opts = append(opts, ometric.WithDescription(desc[0]))
			}
			c, err := o.metricMeter.Int64Counter(counterName, opts...)
			if err != nil {
				return nil, err
			}
			o.counters.Store(counterName, c)
			return c, nil
		})
	}
	if err != nil {
		return err
	}
	if ct == nil {
		return errors.New("nil counter inited")
	}
	ct.(ometric.Int64Counter).Add(ctx, delta)
	return nil
}
