package tsdb_cluster

import (
	"info/internal/domain/concentration"
	"info/internal/infrastructure/repository/tsdb"
)

type ConcentrationReplicaSet struct {
	*ReplicaSet
}

var _ concentration.ReplicaSet = (*ConcentrationReplicaSet)(nil)

func NewConcentrationReplicaSet(replicaSet *ReplicaSet) *ConcentrationReplicaSet {
	return &ConcentrationReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *ConcentrationReplicaSet) WriteRepo() concentration.WriteRepository {
	return tsdb.NewConcentrationRepository(c.ReplicaSet.WriteRepo())
}

func (c *ConcentrationReplicaSet) ReadRepo() concentration.ReadRepository {
	return tsdb.NewConcentrationRepository(c.ReplicaSet.ReadRepo())
}
