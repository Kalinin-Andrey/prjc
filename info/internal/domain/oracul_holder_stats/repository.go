package oracul_holder_stats

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *OraculHolderStats) error
}

type ReadRepository interface {
}
