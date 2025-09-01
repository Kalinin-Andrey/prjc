package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pgx "github.com/jackc/pgx/v5"
	"info/internal/domain"
	"info/internal/domain/price_and_cap"
	"info/internal/pkg/apperror"
	"strconv"
	"strings"
	"time"
)

type PriceAndCapRepository struct {
	*Repository
}

const ()

var _ price_and_cap.WriteRepository = (*PriceAndCapRepository)(nil)
var _ price_and_cap.ReadRepository = (*PriceAndCapRepository)(nil)

func NewPriceAndCapRepository(repository *Repository) *PriceAndCapRepository {
	return &PriceAndCapRepository{
		Repository: repository,
	}
}

const (
	MUpsertPriceAndCap_Limit = 13000 // 5 пар-ра * 13т = 65т ~= max

	price_and_cap_sql_MGet                       = "SELECT currency_id, price, daily_volume, cap, ts FROM cmc.price_and_cap WHERE currency_id = any($1) ORDER BY ts DESC;"
	price_and_cap_sql_Upsert                     = "INSERT INTO cmc.price_and_cap(currency_id, price, daily_volume, cap, ts) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (currency_id, ts) DO UPDATE SET price = EXCLUDED.price, daily_volume = EXCLUDED.daily_volume, cap = EXCLUDED.cap;"
	price_and_cap_sql_MUpsert                    = "INSERT INTO cmc.price_and_cap(currency_id, price, daily_volume, cap, ts) VALUES "
	price_and_cap_sql_MUpsert_OnConflictDoUpdate = " ON CONFLICT (currency_id, ts) DO UPDATE SET price = EXCLUDED.price, daily_volume = EXCLUDED.daily_volume, cap = EXCLUDED.cap;"
)

func (r *PriceAndCapRepository) MGet(ctx context.Context, currencyIDs *[]uint) (price_and_cap.PriceAndCapMap, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PriceAndCapRepository.MGet"

	var entity price_and_cap.PriceAndCap
	res := make(price_and_cap.PriceAndCapMap, len(*currencyIDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, price_and_cap_sql_MGet, *currencyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.Price, &entity.DailyVolume, &entity.Cap, &entity.Ts); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_MGet, err)
		}
		if _, ok := res[entity.CurrencyID]; !ok {
			res[entity.CurrencyID] = make(price_and_cap.PriceAndCapList, 0, defaultCapacityForResult)
		}
		res[entity.CurrencyID] = append(res[entity.CurrencyID], entity)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return res, nil
}

func (r *PriceAndCapRepository) Upsert(ctx context.Context, entity *price_and_cap.PriceAndCap) (err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PriceAndCapRepository.Upsert"
	start := time.Now().UTC()

	if _, err := r.db.Exec(ctx, price_and_cap_sql_Upsert, entity.CurrencyID, entity.Price, entity.DailyVolume, entity.Cap, entity.Ts); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_Upsert, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r PriceAndCapRepository) MUpsertTx(ctx context.Context, tx domain.Tx, entities *[]price_and_cap.PriceAndCap) error {
	if len(*entities) <= MUpsertPriceAndCap_Limit {
		return r.mUpsertTx(ctx, tx, entities)
	}

	lbound := 0
	hbound := MUpsertPriceAndCap_Limit
	for lbound < hbound {
		entitiesItem := (*entities)[lbound:hbound]
		if err := r.mUpsertTx(ctx, tx, &entitiesItem); err != nil {
			return err
		}
		lbound = hbound
		hbound += MUpsertPriceAndCap_Limit
		if hbound > len(*entities) {
			hbound = len(*entities)
		}
	}
	return nil
}

func (r PriceAndCapRepository) mUpsertTx(ctx context.Context, tx domain.Tx, entities *[]price_and_cap.PriceAndCap) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PriceAndCapRepository.mUpsertTx"
	const fields_nb = 5 // при изменении количества полей нужно изменить MUpsertNmDimensions_Limit, чтобы, в результате, кол-во пар-ов не превышало 65т
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(price_and_cap_sql_MUpsert)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ", $" + strconv.Itoa(i*fields_nb+4) + ", $" + strconv.Itoa(i*fields_nb+5) + ")")
		params = append(params, entity.CurrencyID, entity.Price, entity.DailyVolume, entity.Cap, entity.Ts)
	}
	b.WriteString(price_and_cap_sql_MUpsert_OnConflictDoUpdate)
	start := time.Now().UTC()

	_, err := tx.Exec(ctx, b.String(), params...)
	if err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, b.String(), err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
