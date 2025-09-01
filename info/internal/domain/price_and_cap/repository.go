package price_and_cap

import (
	"context"
	"info/internal/domain"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *PriceAndCap) (err error)
	MUpsertTx(ctx context.Context, tx domain.Tx, entities *[]PriceAndCap) error
}

type ReadRepository interface {
	MGet(ctx context.Context, currencyIDs *[]uint) (PriceAndCapMap, error)
}
