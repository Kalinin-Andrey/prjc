package tsdb_cluster

import (
	"info/internal/domain/oracul_daily_balance_stats"
	"info/internal/infrastructure/repository/tsdb"
)

type OraculDailyBalanceStatsReplicaSet struct {
	*ReplicaSet
}

var _ oracul_daily_balance_stats.ReplicaSet = (*OraculDailyBalanceStatsReplicaSet)(nil)

func NewOraculDailyBalanceStatsReplicaSet(replicaSet *ReplicaSet) *OraculDailyBalanceStatsReplicaSet {
	return &OraculDailyBalanceStatsReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *OraculDailyBalanceStatsReplicaSet) WriteRepo() oracul_daily_balance_stats.WriteRepository {
	return tsdb.NewOraculDailyBalanceStatsRepository(c.ReplicaSet.WriteRepo())
}

func (c *OraculDailyBalanceStatsReplicaSet) ReadRepo() oracul_daily_balance_stats.ReadRepository {
	return tsdb.NewOraculDailyBalanceStatsRepository(c.ReplicaSet.ReadRepo())
}
