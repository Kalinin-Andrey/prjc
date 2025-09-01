package currency

import (
	"context"
	"info/internal/domain"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Begin(ctx context.Context) (domain.Tx, error)
	GetImportMaxTimeForUpdateTx(ctx context.Context, tx domain.Tx, currencyIDs *[]uint) (map[uint]ImportMaxTime, error)
	Create(ctx context.Context, entity *Currency) (ID uint, err error)
	MUpsert(ctx context.Context, entities *CurrencyList) error
	Update(ctx context.Context, entity *Currency) error
	Delete(ctx context.Context, ID uint) error
	MCreateImportMaxTime(ctx context.Context, entities *[]ImportMaxTime) error
	MUpsertImportMaxTimeTx(ctx context.Context, tx domain.Tx, entities *[]ImportMaxTime) error
	MUpsertImportMaxTimeMapTx(ctx context.Context, tx domain.Tx, entities map[uint]ImportMaxTime) error
}

type ReadRepository interface {
	Get(ctx context.Context, ID uint) (*Currency, error)
	GetBySlug(ctx context.Context, slug string) (*Currency, error)
	MGet(ctx context.Context, IDs *[]uint) (*CurrencyList, error)
	MGetBySlug(ctx context.Context, slugs *[]string) (*CurrencyList, error)
	GetAll(ctx context.Context) (*CurrencyList, error)
	MGetTokenAddress(ctx context.Context, IDs *[]uint) (*TokenAddressList, error)
}
