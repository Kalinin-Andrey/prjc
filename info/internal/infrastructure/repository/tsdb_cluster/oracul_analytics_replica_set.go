package tsdb_cluster

import (
	"info/internal/domain/oracul_analytics"
	"info/internal/infrastructure/repository/tsdb"
)

type OraculAnalyticsReplicaSet struct {
	*ReplicaSet
}

var _ oracul_analytics.ReplicaSet = (*OraculAnalyticsReplicaSet)(nil)

func NewOraculAnalyticsReplicaSet(replicaSet *ReplicaSet) *OraculAnalyticsReplicaSet {
	return &OraculAnalyticsReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *OraculAnalyticsReplicaSet) WriteRepo() oracul_analytics.WriteRepository {
	return tsdb.NewOraculAnalyticsRepository(c.ReplicaSet.WriteRepo())
}

func (c *OraculAnalyticsReplicaSet) ReadRepo() oracul_analytics.ReadRepository {
	return tsdb.NewOraculAnalyticsRepository(c.ReplicaSet.ReadRepo())
}
