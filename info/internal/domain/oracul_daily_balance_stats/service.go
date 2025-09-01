package oracul_daily_balance_stats

import (
	"context"
)

type Service struct {
	replicaSet ReplicaSet
}

func NewService(replicaSet ReplicaSet) *Service {
	return &Service{
		replicaSet: replicaSet,
	}
}

func (s *Service) MCreate(ctx context.Context, entities *OraculDailyBalanceStatsList) error {
	return s.replicaSet.WriteRepo().MUpsert(ctx, entities)
}
