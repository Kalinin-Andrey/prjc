package tsdb_cluster

import (
	"info/internal/domain/oracul_holder_stats"
	"info/internal/infrastructure/repository/tsdb"
)

type OraculHolderStatsReplicaSet struct {
	*ReplicaSet
}

var _ oracul_holder_stats.ReplicaSet = (*OraculHolderStatsReplicaSet)(nil)

func NewOraculHolderStatsReplicaSet(replicaSet *ReplicaSet) *OraculHolderStatsReplicaSet {
	return &OraculHolderStatsReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *OraculHolderStatsReplicaSet) WriteRepo() oracul_holder_stats.WriteRepository {
	return tsdb.NewOraculHolderStatsRepository(c.ReplicaSet.WriteRepo())
}

func (c *OraculHolderStatsReplicaSet) ReadRepo() oracul_holder_stats.ReadRepository {
	return tsdb.NewOraculHolderStatsRepository(c.ReplicaSet.ReadRepo())
}
