package tsdb

import (
	"context"
	"fmt"
	"info/internal/domain/oracul_daily_balance_stats"
	"info/internal/pkg/apperror"
	"strconv"
	"strings"
	"time"
)

type OraculDailyBalanceStatsRepository struct {
	*Repository
}

var _ oracul_daily_balance_stats.WriteRepository = (*OraculDailyBalanceStatsRepository)(nil)
var _ oracul_daily_balance_stats.ReadRepository = (*OraculDailyBalanceStatsRepository)(nil)

func NewOraculDailyBalanceStatsRepository(repository *Repository) *OraculDailyBalanceStatsRepository {
	return &OraculDailyBalanceStatsRepository{
		Repository: repository,
	}
}

const (
	oracul_daily_balance_stats_MUpsert_Limit = 8000 // 8 пар-ра * 8т <= ~65т ~= max

	oracul_daily_balance_stats_sql_MUpsert                    = "INSERT INTO oracul.daily_balance_stats(currency_id, whales_balance, whales_total_holders, investors_balance, investors_total_holders, retailers_balance, retailers_total_holders, d) VALUES "
	oracul_daily_balance_stats_sql_MUpsert_OnConflictDoUpdate = " ON CONFLICT (currency_id, d) DO UPDATE SET whales_balance = EXCLUDED.whales_balance, whales_total_holders = EXCLUDED.whales_total_holders, investors_balance = EXCLUDED.investors_balance, investors_total_holders = EXCLUDED.investors_total_holders, retailers_balance = EXCLUDED.retailers_balance, retailers_total_holders = EXCLUDED.retailers_total_holders;"
)

func (r *OraculDailyBalanceStatsRepository) MUpsert(ctx context.Context, entities *oracul_daily_balance_stats.OraculDailyBalanceStatsList) error {
	if len(*entities) <= oracul_daily_balance_stats_MUpsert_Limit {
		return r.mUpsert(ctx, entities)
	}

	lbound := 0
	hbound := oracul_daily_balance_stats_MUpsert_Limit
	for lbound < hbound {
		entitiesItem := (*entities)[lbound:hbound]
		if err := r.mUpsert(ctx, &entitiesItem); err != nil {
			return err
		}
		lbound = hbound
		hbound += oracul_daily_balance_stats_MUpsert_Limit
		if hbound > len(*entities) {
			hbound = len(*entities)
		}
	}
	return nil
}

func (r *OraculDailyBalanceStatsRepository) mUpsert(ctx context.Context, entities *oracul_daily_balance_stats.OraculDailyBalanceStatsList) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "OraculDailyBalanceStatsRepository.mUpsert"
	const fields_nb = 8 // при изменении количества полей нужно изменить oracul_daily_balance_stats_MUpsert_Limit, чтобы, в результате, кол-во пар-ов не превышало 65т
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(oracul_daily_balance_stats_sql_MUpsert)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ", $" + strconv.Itoa(i*fields_nb+4) + ", $" + strconv.Itoa(i*fields_nb+5) + ", $" + strconv.Itoa(i*fields_nb+6) + ", $" + strconv.Itoa(i*fields_nb+7) + ", $" + strconv.Itoa(i*fields_nb+8) + ")")
		params = append(params, entity.CurrencyID, entity.WhalesBalance, entity.WhalesTotalHolders, entity.InvestorsBalance, entity.InvestorsTotalHolders, entity.RetailersBalance, entity.RetailersTotalHolders, entity.D)
	}
	b.WriteString(oracul_daily_balance_stats_sql_MUpsert_OnConflictDoUpdate)
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, b.String(), params...)
	if err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, b.String(), err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
