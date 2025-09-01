package tsdb

import (
	"context"
	"fmt"
	"time"

	"info/internal/pkg/apperror"

	"info/internal/domain/oracul_speedometers"
)

type OraculSpeedometersRepository struct {
	*Repository
}

var _ oracul_speedometers.WriteRepository = (*OraculSpeedometersRepository)(nil)
var _ oracul_speedometers.ReadRepository = (*OraculSpeedometersRepository)(nil)

func NewOraculSpeedometersRepository(repository *Repository) *OraculSpeedometersRepository {
	return &OraculSpeedometersRepository{
		Repository: repository,
	}
}

const (
	oracul_speedometers_sql_Upsert = "INSERT INTO oracul.speedometers(currency_id, whales_buy_rate, whales_sell_rate, whales_volume, investors_buy_rate, investors_sell_rate, investors_volume, retailers_buy_rate, retailers_sell_rate, retailers_volume, ts) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (currency_id, ts) DO UPDATE SET whales_buy_rate = EXCLUDED.whales_buy_rate, whales_sell_rate = EXCLUDED.whales_sell_rate, whales_volume = EXCLUDED.whales_volume, investors_buy_rate = EXCLUDED.investors_buy_rate, investors_sell_rate = EXCLUDED.investors_sell_rate, investors_volume = EXCLUDED.investors_volume, retailers_buy_rate = EXCLUDED.retailers_buy_rate, retailers_sell_rate = EXCLUDED.retailers_sell_rate, retailers_volume = EXCLUDED.retailers_volume;"
)

func (r *OraculSpeedometersRepository) Upsert(ctx context.Context, entity *oracul_speedometers.OraculSpeedometers) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "OraculSpeedometersRepository.Upsert"
	start := time.Now().UTC()

	if _, err := r.db.Exec(ctx, oracul_speedometers_sql_Upsert, entity.CurrencyID, entity.WhalesBuyRate, entity.WhalesSellRate, entity.WhalesVolume, entity.InvestorsBuyRate, entity.InvestorsSellRate, entity.InvestorsVolume, entity.RetailersBuyRate, entity.RetailersSellRate, entity.RetailersVolume, entity.Ts); err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, oracul_speedometers_sql_Upsert, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
