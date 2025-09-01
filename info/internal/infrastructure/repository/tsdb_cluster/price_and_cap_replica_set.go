package tsdb_cluster

import (
	"info/internal/domain/price_and_cap"
	"info/internal/infrastructure/repository/tsdb"
)

type PriceAndCapReplicaSet struct {
	*ReplicaSet
}

var _ price_and_cap.ReplicaSet = (*PriceAndCapReplicaSet)(nil)

func NewPriceAndCapReplicaSet(replicaSet *ReplicaSet) *PriceAndCapReplicaSet {
	return &PriceAndCapReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *PriceAndCapReplicaSet) WriteRepo() price_and_cap.WriteRepository {
	return tsdb.NewPriceAndCapRepository(c.ReplicaSet.WriteRepo())
}

func (c *PriceAndCapReplicaSet) ReadRepo() price_and_cap.ReadRepository {
	return tsdb.NewPriceAndCapRepository(c.ReplicaSet.ReadRepo())
}
