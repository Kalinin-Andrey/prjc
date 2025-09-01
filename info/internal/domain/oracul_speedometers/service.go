package oracul_speedometers

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

func (s *Service) Create(ctx context.Context, entity *OraculSpeedometers) error {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}
