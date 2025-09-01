package tsdb

import (
	"context"
	"fmt"
	"info/internal/domain/oracul_analytics"
	"info/internal/pkg/apperror"
	"time"
)

type OraculAnalyticsRepository struct {
	*Repository
}

var _ oracul_analytics.WriteRepository = (*OraculAnalyticsRepository)(nil)
var _ oracul_analytics.ReadRepository = (*OraculAnalyticsRepository)(nil)

func NewOraculAnalyticsRepository(repository *Repository) *OraculAnalyticsRepository {
	return &OraculAnalyticsRepository{
		Repository: repository,
	}
}

const (
	oracul_analytics_sql_Upsert = "INSERT INTO oracul.analytics(currency_id, whales_concentration, worm_index, growth_fuel, ts) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (currency_id, ts) DO UPDATE SET whales_concentration = EXCLUDED.whales_concentration, worm_index = EXCLUDED.worm_index, growth_fuel = EXCLUDED.growth_fuel;"
)

func (r *OraculAnalyticsRepository) Upsert(ctx context.Context, entity *oracul_analytics.OraculAnalytics) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "OraculAnalyticsRepository.Upsert"
	start := time.Now().UTC()

	if _, err := r.db.Exec(ctx, oracul_analytics_sql_Upsert, entity.CurrencyID, entity.WhalesConcentration, entity.WormIndex, entity.GrowthFuel, entity.Ts); err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, oracul_analytics_sql_Upsert, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
