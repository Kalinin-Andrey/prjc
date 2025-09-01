package oracul_speedometers

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *OraculSpeedometers) error
}

type ReadRepository interface {
}
