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
	"info/internal/domain/currency"
)

type CurrencyRepository struct {
	*Repository
}

const ()

var _ currency.WriteRepository = (*CurrencyRepository)(nil)
var _ currency.ReadRepository = (*CurrencyRepository)(nil)

func NewCurrencyRepository(repository *Repository) *CurrencyRepository {
	return &CurrencyRepository{
		Repository: repository,
	}
}

const (
	currency_sql_Get                       = "SELECT id, symbol, slug, name, is_for_observing, circulating_supply, self_reported_circulating_supply, total_supply, max_supply, latest_price, cmc_rank, date_added, platform FROM cmc.currency WHERE id = $1;"
	currency_sql_GetBySlug                 = "SELECT id, symbol, slug, name, is_for_observing, circulating_supply, self_reported_circulating_supply, total_supply, max_supply, latest_price, cmc_rank, date_added, platform FROM cmc.currency WHERE slug = $1;"
	currency_sql_GetImportMaxTimeForUpdate = "SELECT currency_id, price_and_cap, concentration FROM cmc.import_max_time WHERE currency_id = ANY($1) FOR UPDATE;"
	currency_sql_MGet                      = "SELECT id, symbol, slug, name, is_for_observing, circulating_supply, self_reported_circulating_supply, total_supply, max_supply, latest_price, cmc_rank, date_added, platform FROM cmc.currency WHERE id = any($1);"
	currency_sql_MGetTokenAddress          = "SELECT currency_id, blockchain, address FROM cmc.token_address WHERE currency_id = any($1);"
	currency_sql_MGetBySlug                = "SELECT id, symbol, slug, name, is_for_observing, circulating_supply, self_reported_circulating_supply, total_supply, max_supply, latest_price, cmc_rank, date_added, platform FROM cmc.currency WHERE slug = any($1);"
	currency_sql_GetAll                    = "SELECT id, symbol, slug, name, is_for_observing, circulating_supply, self_reported_circulating_supply, total_supply, max_supply, latest_price, cmc_rank, date_added, platform FROM cmc.currency WHERE is_for_observing = TRUE;"
	currency_sql_Create                    = "INSERT INTO cmc.currency(id, symbol, slug, name, is_for_observing) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING id;"
	currency_sql_MCreate                   = "INSERT INTO cmc.currency(id, symbol, slug, name, is_for_observing, circulating_supply, self_reported_circulating_supply, total_supply, max_supply, latest_price, cmc_rank, date_added, platform) VALUES "
	currency_sql_Create_OnConflictDoUpdate = " ON CONFLICT (id) DO UPDATE SET symbol = EXCLUDED.symbol, slug = EXCLUDED.slug, name = EXCLUDED.name, is_for_observing = EXCLUDED.is_for_observing, circulating_supply = EXCLUDED.circulating_supply, self_reported_circulating_supply = EXCLUDED.self_reported_circulating_supply, total_supply = EXCLUDED.total_supply, max_supply = EXCLUDED.max_supply, latest_price = EXCLUDED.latest_price, cmc_rank = EXCLUDED.cmc_rank, date_added = EXCLUDED.date_added, platform = EXCLUDED.platform;"
	currency_sql_Update                    = "UPDATE cmc.currency SET symbol = $2, slug = $3, name = $4, is_for_observing = $5 WHERE id = $1;"
	currency_sql_Delete                    = "DELETE FROM cmc.currency WHERE id = $1;"

	import_max_time_sql_MCreate                    = "INSERT INTO cmc.import_max_time(currency_id, price_and_cap, concentration) VALUES "
	import_max_time_sql_MCreate_OnConflictDoUpdate = " ON CONFLICT (currency_id) DO UPDATE SET price_and_cap = EXCLUDED.price_and_cap, concentration = EXCLUDED.concentration;"
)

func (r *CurrencyRepository) Get(ctx context.Context, ID uint) (*currency.Currency, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.Get"
	start := time.Now().UTC()

	entity := &currency.Currency{}
	if err := r.db.QueryRow(ctx, currency_sql_Get, ID).Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving, &entity.CirculatingSupply, &entity.SelfReportedCirculatingSupply, &entity.TotalSupply, &entity.MaxSupply, &entity.LatestPrice, &entity.CmcRank, &entity.AddedAt, &entity.Platform); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Get, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *CurrencyRepository) GetBySlug(ctx context.Context, slug string) (*currency.Currency, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetBySlug"
	start := time.Now().UTC()

	entity := &currency.Currency{}
	if err := r.db.QueryRow(ctx, currency_sql_GetBySlug, slug).Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving, &entity.CirculatingSupply, &entity.SelfReportedCirculatingSupply, &entity.TotalSupply, &entity.MaxSupply, &entity.LatestPrice, &entity.CmcRank, &entity.AddedAt, &entity.Platform); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetBySlug, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *CurrencyRepository) GetImportMaxTimeForUpdateTx(ctx context.Context, tx domain.Tx, currencyIDs *[]uint) (map[uint]currency.ImportMaxTime, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetImportMaxTimeForUpdateTx"

	var entity currency.ImportMaxTime
	res := make(map[uint]currency.ImportMaxTime, len(*currencyIDs))

	start := time.Now().UTC()
	rows, err := tx.Query(ctx, currency_sql_GetImportMaxTimeForUpdate, *currencyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetImportMaxTimeForUpdate, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.PriceAndCap, &entity.Concentration); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetImportMaxTimeForUpdate, err)
		}
		res[entity.CurrencyID] = entity
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return res, nil
}

func (r *CurrencyRepository) MGet(ctx context.Context, IDs *[]uint) (*currency.CurrencyList, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.MGet"

	var entity currency.Currency
	res := make(currency.CurrencyList, 0, len(*IDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_MGet, *IDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving, &entity.CirculatingSupply, &entity.SelfReportedCirculatingSupply, &entity.TotalSupply, &entity.MaxSupply, &entity.LatestPrice, &entity.CmcRank, &entity.AddedAt, &entity.Platform); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGet, err)
		}
		res = append(res, entity)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r *CurrencyRepository) MGetBySlug(ctx context.Context, slugs *[]string) (*currency.CurrencyList, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.MGetBySlug"

	var entity currency.Currency
	res := make(currency.CurrencyList, 0, len(*slugs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_MGetBySlug, *slugs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGetBySlug, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving, &entity.CirculatingSupply, &entity.SelfReportedCirculatingSupply, &entity.TotalSupply, &entity.MaxSupply, &entity.LatestPrice, &entity.CmcRank, &entity.AddedAt, &entity.Platform); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGetBySlug, err)
		}
		res = append(res, entity)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r *CurrencyRepository) GetAll(ctx context.Context) (*currency.CurrencyList, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetAll"

	var entity currency.Currency
	res := make(currency.CurrencyList, 0, defaultCapacityForResult)

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_GetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetAll, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving, &entity.CirculatingSupply, &entity.SelfReportedCirculatingSupply, &entity.TotalSupply, &entity.MaxSupply, &entity.LatestPrice, &entity.CmcRank, &entity.AddedAt, &entity.Platform); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetAll, err)
		}
		res = append(res, entity)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r *CurrencyRepository) Create(ctx context.Context, entity *currency.Currency) (ID uint, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "CurrencyRepository.Create"
	start := time.Now().UTC()

	if err := r.db.QueryRow(ctx, currency_sql_Create, entity.ID, entity.Symbol, entity.Slug, entity.Name, entity.IsForObserving).Scan(&ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return 0, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return 0, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Create, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return ID, nil
}

func (r *CurrencyRepository) Update(ctx context.Context, entity *currency.Currency) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.Update"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, currency_sql_Update, entity.ID, entity.Symbol, entity.Slug, entity.Name, entity.IsForObserving)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Update, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *CurrencyRepository) Delete(ctx context.Context, ID uint) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.Delete"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, currency_sql_Delete, ID)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Delete, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r CurrencyRepository) MCreateImportMaxTime(ctx context.Context, entities *[]currency.ImportMaxTime) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PriceAndCapRepository.MCreateImportMaxTime"
	const fields_nb = 3 // при изменении количества полей нужно изменить MUpsertNmDimensions_Limit, чтобы, в результате, кол-во пар-ов не превышало 65т
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(import_max_time_sql_MCreate)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ")")
		params = append(params, entity.CurrencyID, entity.PriceAndCap, entity.Concentration)
	}
	b.WriteString(sql_OnConflictDoNothing)
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

func (r CurrencyRepository) MUpsertImportMaxTimeMapTx(ctx context.Context, tx domain.Tx, entities map[uint]currency.ImportMaxTime) error {
	list := make([]currency.ImportMaxTime, 0, len(entities))
	var item currency.ImportMaxTime
	for _, item = range entities {
		list = append(list, item)
	}
	return r.MUpsertImportMaxTimeTx(ctx, tx, &list)
}

func (r CurrencyRepository) MUpsertImportMaxTimeTx(ctx context.Context, tx domain.Tx, entities *[]currency.ImportMaxTime) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PriceAndCapRepository.MUpsertImportMaxTime"
	const fields_nb = 3 // при изменении количества полей нужно изменить MUpsertNmDimensions_Limit, чтобы, в результате, кол-во пар-ов не превышало 65т
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(import_max_time_sql_MCreate)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ")")
		params = append(params, entity.CurrencyID, entity.PriceAndCap, entity.Concentration)
	}
	b.WriteString(import_max_time_sql_MCreate_OnConflictDoUpdate)
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

func (r CurrencyRepository) MUpsert(ctx context.Context, entities *currency.CurrencyList) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "WarehouseRepository.MUpsert"
	const fields_nb = 13
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(currency_sql_MCreate)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ", $" + strconv.Itoa(i*fields_nb+4) + ", $" + strconv.Itoa(i*fields_nb+5) + ", $" + strconv.Itoa(i*fields_nb+6) + ", $" + strconv.Itoa(i*fields_nb+7) + ", $" + strconv.Itoa(i*fields_nb+8) + ", $" + strconv.Itoa(i*fields_nb+9) + ", $" + strconv.Itoa(i*fields_nb+10) + ", $" + strconv.Itoa(i*fields_nb+11) + ", $" + strconv.Itoa(i*fields_nb+12) + ", $" + strconv.Itoa(i*fields_nb+13) + ")")
		params = append(params, entity.ID, entity.Symbol, entity.Slug, entity.Name, entity.IsForObserving, entity.CirculatingSupply, entity.SelfReportedCirculatingSupply, entity.TotalSupply, entity.MaxSupply, entity.LatestPrice, entity.CmcRank, entity.AddedAt, entity.Platform)
	}
	b.WriteString(currency_sql_Create_OnConflictDoUpdate)
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

func (r *CurrencyRepository) MGetTokenAddress(ctx context.Context, IDs *[]uint) (*currency.TokenAddressList, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.MGetTokenAddress"

	var entity currency.TokenAddress
	res := make(currency.TokenAddressList, 0, len(*IDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_MGetTokenAddress, *IDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGetTokenAddress, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.Blockchain, &entity.Address); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGetTokenAddress, err)
		}
		res = append(res, entity)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

/*
ЗАПРОСЫ


Концентрация у китов, мин/макс цена (для таблички)

with d as (
select currency_id, max(d) as last_d
from cmc.concentration
group by currency_id),
whales as (
select currency_id, d, whales * 100 / (whales + investors + retail) as whales_prc
from concentration),
ts as (
select currency_id, max(ts) as last_ts
from cmc.price_and_cap
group by currency_id),
min_max as (
select currency_id, min(price) as min_price, max(price) as max_price
from cmc.price_and_cap
group by currency_id),
price as (
select currency_id, price, ts
from cmc.price_and_cap
)
select c.id, c.symbol, round(cast(whales.whales_prc AS numeric), 2) as whales_prc, price.price, min_max.min_price, min_max.max_price,
round(cast(min_max.max_price/price.price AS numeric)) as to_ath, round(cast(price.price/min_max.min_price AS numeric)) as from_atl
from currency c
inner join ts on c.id = ts.currency_id
inner join price on ts.currency_id = price.currency_id and ts.last_ts = price.ts
inner join d on c.id = d.currency_id
inner join whales on d.currency_id = whales.currency_id and d.last_d = whales.d
inner join min_max on c.id = min_max.currency_id;


Изменения концентрации у китов [1неделя, 2недели, 1месяц, 2месяца]

with m2 as (
select currency_id, max(d) as d
from cmc.concentration
where d <= (now() - interval '2 month')::date
group by currency_id),
m1 as (
select currency_id, max(d) as d
from cmc.concentration
where d <= (now() - interval '1 month')::date
group by currency_id),
w2 as (
select currency_id, max(d) as d
from cmc.concentration
where d <= (now() - interval '2 week')::date
group by currency_id),
w1 as (
select currency_id, max(d) as d
from cmc.concentration
where d <= (now() - interval '1 week')::date
group by currency_id),
n as (
select currency_id, max(d) as d
from cmc.concentration
where d <= (now() - interval '2 day')::date		-- из-за глюков в данных
group by currency_id)
select c.id, c.symbol,
round(cast(((cn.whales - cw1.whales) * 100)/cw1.whales AS numeric), 4) as week1, round(cast(((cn.whales - cw2.whales) * 100)/cw2.whales AS numeric), 4) as week2,
round(cast(((cn.whales - cm1.whales) * 100)/cm1.whales AS numeric), 4) as month1, round(cast(((cn.whales - cm2.whales) * 100)/cm2.whales AS numeric), 4) as month2
from cmc.currency c
inner join m2 on c.id = m2.currency_id
inner join m1 on c.id = m1.currency_id
inner join w2 on c.id = w2.currency_id
inner join w1 on c.id = w1.currency_id
inner join n on c.id = n.currency_id
inner join cmc.concentration cm2 on m2.currency_id = cm2.currency_id and m2.d = cm2.d
inner join cmc.concentration cm1 on m1.currency_id = cm1.currency_id and m1.d = cm1.d
inner join cmc.concentration cw2 on w2.currency_id = cw2.currency_id and w2.d = cw2.d
inner join cmc.concentration cw1 on w1.currency_id = cw1.currency_id and w1.d = cw1.d
inner join cmc.concentration cn on n.currency_id = cn.currency_id and n.d = cn.d
order by month2, month1, week2, week1;

Запросы для графиков

Изменения доли китов
WITH val AS (
	SELECT c.currency_id, c.d, c.whales
	FROM
		cmc.concentration AS c
		INNER JOIN (
			SELECT currency_id, max(d) AS d
			FROM cmc.concentration
			WHERE d <= (now() - interval '$period')::date
			GROUP BY currency_id
		) AS dt ON c.currency_id = dt.currency_id AND c.d = dt.d
)
SELECT v.d AS time, c.symbol AS text, Round(Cast(((v.whales - val.whales) * 100)/val.whales AS numeric), 4) AS w
FROM
	cmc.currency c
	INNER JOIN val ON c.id = val.currency_id
	INNER JOIN cmc.concentration v ON v.currency_id = val.currency_id AND v.d >= val.d AND v.d <= (now() - interval '2 day')::date
WHERE c.is_for_observing = true
GROUP BY text, time, w
ORDER BY time;

Доля китов
SELECT
  $__timeGroup(conc.d,'1d',previous),
  c.symbol as text,
  (conc.whales * 100)/(conc.whales + conc.investors + conc.retail) as w
FROM
  cmc.concentration AS conc
  INNER JOIN cmc.currency AS c ON c.id = conc.currency_id
WHERE c.is_for_observing = true AND conc.whales IS NOT Null AND conc.whales > 0 AND conc.d <= (now() - interval '2 day')::date AND $__timeFilter(conc.d)
GROUP BY text, time, w
ORDER BY time;

Изменения дневного объёма
WITH val AS (
	SELECT c.currency_id, c.ts, c.daily_volume
	FROM
		cmc.price_and_cap AS c
		INNER JOIN (
			SELECT currency_id, max(ts) AS ts
			FROM cmc.price_and_cap
			WHERE ts <= (now() - interval '$period')::date
			GROUP BY currency_id
		) AS dt ON c.currency_id = dt.currency_id AND c.ts = dt.ts
)
SELECT v.ts AS time, c.symbol AS text, Round(Cast(((v.daily_volume - val.daily_volume) * 100)/val.daily_volume AS numeric), 4) AS w
FROM
	cmc.currency c
	INNER JOIN val ON c.id = val.currency_id
	INNER JOIN cmc.price_and_cap v ON v.currency_id = val.currency_id AND v.ts >= val.ts AND v.ts <= (now() - interval '2 day')::date
WHERE c.is_for_observing = true
GROUP BY text, time, w
ORDER BY time;

Дневной объём торгов
SELECT
  $__timeGroup(cap.ts,'1d',previous),
  c.symbol as text,
  cap.daily_volume as w
FROM
  cmc.price_and_cap AS cap
  INNER JOIN cmc.currency AS c ON c.id = cap.currency_id
WHERE c.is_for_observing = true AND cap.daily_volume IS NOT Null AND cap.daily_volume > 0 AND cap.ts <= (now() - interval '2 day')::date AND $__timeFilter(cap.ts)
GROUP BY text, time, w
ORDER BY time;



Сводная таблица
SELECT c.symbol, w.whales_prc, ((GREATEST(c.circulating_supply, c.self_reported_circulating_supply) * c.latest_price)/1000000)::integer AS cap, ((coalesce(c.max_supply, c.total_supply) * c.latest_price)/1000000)::integer AS fdv, w.to_ath, w.from_atl, d.bonus, round(cast(coalesce(pit.crypto_holdings, 0) AS numeric), 2) AS crypto_holdings, round(cast(coalesce(pit.pl_percent_value, 0) AS numeric), 2) AS pl_percent_value
FROM
	currency AS c
	LEFT JOIN whales_prc_and_min_max_price AS w ON c.id = w.id
	LEFT JOIN cm1 AS d ON c.id = d.id
	LEFT JOIN portfolio_item AS pit ON pit.portfolio_source_id = '6651f947db928013879d191c' AND c.id = pit.currency_id
WHERE c.is_for_observing = true
ORDER BY c.cmc_rank;

Сводная таблица V2
select c.symbol, w.whales_prc, ((GREATEST(c.circulating_supply, c.self_reported_circulating_supply) * c.latest_price)/1000000)::integer as cap, ((coalesce(c.max_supply, c.total_supply) * c.latest_price)/1000000)::integer as fdv, w.to_ath, w.from_atl, oa.whales_concentration, d.bonus, round(cast(coalesce(pit.crypto_holdings, 0) AS numeric), 2) as crypto_holdings, round(cast(coalesce(pit.pl_percent_value, 0) AS numeric), 2) as pl_percent_value
from currency as c
left join whales_prc_and_min_max_price as w on c.id = w.id
left join cw2 as d on c.id = d.id
left join portfolio_item as pit on pit.portfolio_source_id = '6651f947db928013879d191c' and c.id = pit.currency_id
left join (
	select oa.currency_id, oa.whales_concentration
	from oracul.analytics as oa
	inner join (
		select currency_id, max(ts) as ts
		from oracul.analytics
		group by currency_id
	) as d on d.currency_id = oa.currency_id and d.ts = oa.ts
) as oa on c.id = oa.currency_id
where c.is_for_observing = true
order by c.cmc_rank;

Сводная таблица V3
select c.symbol, w.whales_prc, ((GREATEST(c.circulating_supply, c.self_reported_circulating_supply) * c.latest_price)/1000000)::integer as cap, ((coalesce(c.max_supply, c.total_supply) * c.latest_price)/1000000)::integer as fdv, w.to_ath, w.from_atl, oa.whales_concentration, round(cast(coalesce(pit.buy_avg_price, 0) AS numeric), 2) as buy_avg_price, round(cast(coalesce(pit.current_price, 0) AS numeric), 2) as current_price, round(cast(coalesce(pit.pl_percent_value, 0) AS numeric), 2)*100 as pl_percent_value, round(cast(coalesce(pit.holdings_percent, 0) AS numeric), 2)*100 as holdings_percent, round(cast(coalesce(pit.total_buy_spent, 0) AS numeric), 2) as total_buy_spent, round(cast(coalesce(pit.crypto_holdings, 0) AS numeric), 2) as crypto_holdings
from currency as c
left join whales_prc_and_min_max_price as w on c.id = w.id
left join cw2 as d on c.id = d.id
left join portfolio_item as pit on pit.portfolio_source_id = '6651f947db928013879d191c' and c.id = pit.currency_id
left join (
	select oa.currency_id, oa.whales_concentration
	from oracul.analytics as oa
	inner join (
		select currency_id, max(ts) as ts
		from oracul.analytics
		group by currency_id
	) as d on d.currency_id = oa.currency_id and d.ts = oa.ts
) as oa on c.id = oa.currency_id
where c.is_for_observing = true
order by c.cmc_rank;

Сводная таблица V4
select c.symbol, w.whales_prc, ((GREATEST(c.circulating_supply, c.self_reported_circulating_supply) * c.latest_price)/1000000)::integer as cap, ((coalesce(c.max_supply, c.total_supply) * c.latest_price)/1000000)::integer as fdv, w.to_ath, w.from_atl, oa.whales_concentration, round(cast(coalesce(pit.buy_avg_price, 0) AS numeric), 2) as buy_avg_price, round(cast(coalesce(pit.current_price, 0) AS numeric), 2) as current_price, round(cast(coalesce(pit.pl_percent_value, 0) * 100 AS numeric), 2) as pl_percent_value, round(cast(coalesce(pit.holdings_percent, 0) * 100 AS numeric), 2) as holdings_percent, round(cast(coalesce(pit.total_buy_spent, 0) AS numeric), 2) as total_buy_spent, round(cast(coalesce(pit.crypto_holdings, 0) AS numeric), 2) as crypto_holdings
from currency as c
left join whales_prc_and_min_max_price as w on c.id = w.id
left join cw2 as d on c.id = d.id
left join portfolio_item as pit on pit.portfolio_source_id = '6651f947db928013879d191c' and c.id = pit.currency_id
left join (
	select oa.currency_id, oa.whales_concentration
	from oracul.analytics as oa
	inner join (
		select currency_id, max(ts) as ts
		from oracul.analytics
		group by currency_id
	) as d on d.currency_id = oa.currency_id and d.ts = oa.ts
) as oa on c.id = oa.currency_id
where c.is_for_observing = true
order by c.cmc_rank;


select oa.currency_id, oa.whales_concentration
from oracul.analytics as oa
inner join (
	select currency_id, max(ts) as ts
	from oracul.analytics
	group by currency_id
) as d on d.currency_id = oa.currency_id and d.ts = oa.ts;

//	[Oracul] Изменение за период
with val as (
	select c.currency_id, c.d, c.whales_balance
	from oracul.daily_balance_stats as c
	inner join (
		select currency_id, max(d) as d
		from oracul.daily_balance_stats
		where d <= (now() - interval '$period')::date
		group by currency_id
	) as dt on c.currency_id = dt.currency_id and c.d = dt.d
)
select v.d as time, c.symbol as text, round(cast(((v.whales_balance - val.whales_balance) * 100)/val.whales_balance AS numeric), 4) as w
from cmc.currency c
inner join val on c.id = val.currency_id
inner join oracul.daily_balance_stats v on v.currency_id = val.currency_id and v.d >= val.d and v.d <= (now() - interval '2 day')::date
where c.is_for_observing = true AND val.whales_balance > 0
GROUP BY text, time, w
ORDER BY time;


//	[Oracul] Доля китов (%)
SELECT
  $__timeGroup(conc.d,'1d',previous),
  c.symbol as text,
  (conc.whales_balance * 100)/(conc.whales_balance + conc.investors_balance + conc.retailers_balance) as w
FROM
  oracul.daily_balance_stats AS conc
  INNER JOIN cmc.currency AS c ON c.id = conc.currency_id
WHERE c.is_for_observing = true AND conc.whales_balance IS NOT Null AND conc.whales_balance > 0 AND conc.d <= (now() - interval '1 day')::date AND $__timeFilter(conc.d)
GROUP BY text, time, w
ORDER BY time;


*/
