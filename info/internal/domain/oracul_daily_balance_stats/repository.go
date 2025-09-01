package oracul_daily_balance_stats

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
	MUpsert(ctx context.Context, entities *OraculDailyBalanceStatsList) error
}

type ReadRepository interface {
}
