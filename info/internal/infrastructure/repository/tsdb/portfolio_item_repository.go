package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"info/internal/domain/portfolio_item"
	"info/internal/pkg/apperror"
	"strconv"
	"strings"
	"time"
)

type PortfolioItemRepository struct {
	*Repository
}

const ()

var _ portfolio_item.WriteRepository = (*PortfolioItemRepository)(nil)
var _ portfolio_item.ReadRepository = (*PortfolioItemRepository)(nil)

func NewPortfolioItemRepository(repository *Repository) *PortfolioItemRepository {
	return &PortfolioItemRepository{
		Repository: repository,
	}
}

const (
	portfolio_item_sql_MGet                      = "SELECT portfolio_source_id, currency_id, amount, current_price, crypto_holdings, holdings_percent, buy_avg_price, pl_percent_value, pl_value, total_buy_spent, updated_at FROM cmc.portfolio_item WHERE portfolio_source_id = $1;"
	portfolio_item_sql_MCreate                   = "INSERT INTO cmc.portfolio_item(portfolio_source_id, currency_id, amount, current_price, crypto_holdings, holdings_percent, buy_avg_price, pl_percent_value, pl_value, total_buy_spent, updated_at) VALUES "
	portfolio_item_sql_Create_OnConflictDoUpdate = " ON CONFLICT (portfolio_source_id, currency_id) DO UPDATE SET amount = EXCLUDED.amount, current_price = EXCLUDED.current_price, crypto_holdings = EXCLUDED.crypto_holdings, holdings_percent = EXCLUDED.holdings_percent, buy_avg_price = EXCLUDED.buy_avg_price, pl_percent_value = EXCLUDED.pl_percent_value, pl_value = EXCLUDED.pl_value, total_buy_spent = EXCLUDED.total_buy_spent, updated_at = EXCLUDED.updated_at;"
)

func (r *PortfolioItemRepository) MGetByPortfolioSourceId(ctx context.Context, portfolioSourceId uint) (*portfolio_item.PortfolioItemMap, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PortfolioItemRepository.MGetByPortfolioSourceId"

	var entity portfolio_item.PortfolioItem
	res := make(portfolio_item.PortfolioItemMap, defaultCapacityForResult)

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, portfolio_item_sql_MGet, portfolioSourceId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, portfolio_item_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.PortfolioSourceID, &entity.CurrencyID, &entity.Amount, &entity.CurrentPrice, &entity.CryptoHoldings, &entity.HoldingsPercent, &entity.BuyAvgPrice, &entity.PlPercentValue, &entity.PlValue, &entity.TotalBuySpent, &entity.UpdatedAt); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, portfolio_item_sql_MGet, err)
		}
		res[entity.CurrencyID] = entity
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r PortfolioItemRepository) MUpsert(ctx context.Context, entities *portfolio_item.PortfolioItemList) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PortfolioItemRepository.MUpsert"
	const fields_nb = 11
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(portfolio_item_sql_MCreate)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ", $" + strconv.Itoa(i*fields_nb+4) + ", $" + strconv.Itoa(i*fields_nb+5) + ", $" + strconv.Itoa(i*fields_nb+6) + ", $" + strconv.Itoa(i*fields_nb+7) + ", $" + strconv.Itoa(i*fields_nb+8) + ", $" + strconv.Itoa(i*fields_nb+9) + ", $" + strconv.Itoa(i*fields_nb+10) + ", $" + strconv.Itoa(i*fields_nb+11) + ")")
		params = append(params, entity.PortfolioSourceID, entity.CurrencyID, entity.Amount, entity.CurrentPrice, entity.CryptoHoldings, entity.HoldingsPercent, entity.BuyAvgPrice, entity.PlPercentValue, entity.PlValue, entity.TotalBuySpent, entity.UpdatedAt)
	}
	b.WriteString(portfolio_item_sql_Create_OnConflictDoUpdate)
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
