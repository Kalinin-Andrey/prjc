package tsdb

import (
	"context"
	"fmt"
	"info/internal/domain/oracul_holder_stats"
	"info/internal/pkg/apperror"
	"time"
)

type OraculHolderStatsRepository struct {
	*Repository
}

var _ oracul_holder_stats.WriteRepository = (*OraculHolderStatsRepository)(nil)
var _ oracul_holder_stats.ReadRepository = (*OraculHolderStatsRepository)(nil)

func NewOraculHolderStatsRepository(repository *Repository) *OraculHolderStatsRepository {
	return &OraculHolderStatsRepository{
		Repository: repository,
	}
}

const (
	oracul_holder_stats_sql_Upsert = "INSERT INTO oracul.holder_stats(currency_id, whales_volume, whales_total_holders, investors_volume, investors_total_holders, retailers_volume, retailers_total_holders, ts) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (currency_id, ts) DO UPDATE SET whales_volume = EXCLUDED.whales_volume, whales_total_holders = EXCLUDED.whales_total_holders, investors_volume = EXCLUDED.investors_volume, investors_total_holders = EXCLUDED.investors_total_holders, retailers_volume = EXCLUDED.retailers_volume, retailers_total_holders = EXCLUDED.retailers_total_holders;"
)

func (r *OraculHolderStatsRepository) Upsert(ctx context.Context, entity *oracul_holder_stats.OraculHolderStats) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "OraculHolderStatsRepository.Upsert"
	start := time.Now().UTC()

	if _, err := r.db.Exec(ctx, oracul_holder_stats_sql_Upsert, entity.CurrencyID, entity.WhalesVolume, entity.WhalesTotalHolders, entity.InvestorsVolume, entity.InvestorsTotalHolders, entity.RetailersVolume, entity.RetailersTotalHolders, entity.Ts); err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, oracul_holder_stats_sql_Upsert, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
