package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	pgx "github.com/jackc/pgx/v5"

	"info/internal/pkg/apperror"

	"info/internal/domain"
	"info/internal/domain/concentration"
)

type ConcentrationRepository struct {
	*Repository
}

const ()

var _ concentration.WriteRepository = (*ConcentrationRepository)(nil)
var _ concentration.ReadRepository = (*ConcentrationRepository)(nil)

func NewConcentrationRepository(repository *Repository) *ConcentrationRepository {
	return &ConcentrationRepository{
		Repository: repository,
	}
}

const (
	MUpsertConcentration_Limit = 11000 // 6 пар-ра * 13т = 65т ~= max

	concentration_sql_MGet                       = "SELECT currency_id, whales, investors, retail, d FROM cmc.concentration WHERE currency_id = any($1) ORDER BY d DESC;"
	concentration_sql_Upsert                     = "INSERT INTO cmc.concentration(currency_id, whales, investors, retail, d) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (currency_id, d) DO UPDATE SET whales = EXCLUDED.whales, investors = EXCLUDED.investors, retail = EXCLUDED.retail;"
	concentration_sql_MUpsert                    = "INSERT INTO cmc.concentration(currency_id, whales, investors, retail, d) VALUES "
	concentration_sql_MUpsert_OnConflictDoUpdate = " ON CONFLICT (currency_id, d) DO UPDATE SET whales = EXCLUDED.whales, investors = EXCLUDED.investors, retail = EXCLUDED.retail;"
)

func (r *ConcentrationRepository) MGet(ctx context.Context, currencyIDs *[]uint) (concentration.ConcentrationMap, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "ConcentrationRepository.MGet"

	var entity concentration.Concentration
	res := make(concentration.ConcentrationMap, len(*currencyIDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, concentration_sql_MGet, *currencyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.Whales, &entity.Investors, &entity.Retail, &entity.D); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_MGet, err)
		}
		if _, ok := res[entity.CurrencyID]; !ok {
			res[entity.CurrencyID] = make(concentration.ConcentrationList, 0, defaultCapacityForResult)
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

func (r *ConcentrationRepository) Upsert(ctx context.Context, entity *concentration.Concentration) (err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "ConcentrationRepository.Upsert"
	start := time.Now().UTC()

	if _, err := r.db.Exec(ctx, concentration_sql_Upsert, entity.CurrencyID, entity.Whales, entity.Investors, entity.Retail, entity.D); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_Upsert, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *ConcentrationRepository) MUpsertTx(ctx context.Context, tx domain.Tx, entities *[]concentration.Concentration) error {
	if len(*entities) <= MUpsertConcentration_Limit {
		return r.mUpsertTx(ctx, tx, entities)
	}

	lbound := 0
	hbound := MUpsertConcentration_Limit
	for lbound < hbound {
		entitiesItem := (*entities)[lbound:hbound]
		if err := r.mUpsertTx(ctx, tx, &entitiesItem); err != nil {
			return err
		}
		lbound = hbound
		hbound += MUpsertConcentration_Limit
		if hbound > len(*entities) {
			hbound = len(*entities)
		}
	}
	return nil
}

func (r *ConcentrationRepository) mUpsertTx(ctx context.Context, tx domain.Tx, entities *[]concentration.Concentration) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "ConcentrationRepository.mUpsertTx"
	const fields_nb = 5 // при изменении количества полей нужно изменить MUpsertNmDimensions_Limit, чтобы, в результате, кол-во пар-ов не превышало 65т
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(concentration_sql_MUpsert)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ", $" + strconv.Itoa(i*fields_nb+4) + ", $" + strconv.Itoa(i*fields_nb+5) + ")")
		params = append(params, entity.CurrencyID, entity.Whales, entity.Investors, entity.Retail, entity.D)
	}
	b.WriteString(concentration_sql_MUpsert_OnConflictDoUpdate)
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
