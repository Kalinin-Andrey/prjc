package concentration

import (
	"context"
	"info/internal/domain"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *Concentration) (err error)
	MUpsertTx(ctx context.Context, tx domain.Tx, entities *[]Concentration) error
}

type ReadRepository interface {
	MGet(ctx context.Context, currencyIDs *[]uint) (ConcentrationMap, error)
}
