package tsdb_cluster

import (
	"info/internal/domain/currency"
	"info/internal/infrastructure/repository/tsdb"
)

type CurrencyReplicaSet struct {
	*ReplicaSet
}

var _ currency.ReplicaSet = (*CurrencyReplicaSet)(nil)

func NewCurrencyReplicaSet(replicaSet *ReplicaSet) *CurrencyReplicaSet {
	return &CurrencyReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *CurrencyReplicaSet) WriteRepo() currency.WriteRepository {
	return tsdb.NewCurrencyRepository(c.ReplicaSet.WriteRepo())
}

func (c *CurrencyReplicaSet) ReadRepo() currency.ReadRepository {
	return tsdb.NewCurrencyRepository(c.ReplicaSet.ReadRepo())
}
