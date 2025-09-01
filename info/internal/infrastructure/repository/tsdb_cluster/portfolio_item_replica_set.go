package tsdb_cluster

import (
	"info/internal/domain/portfolio_item"
	"info/internal/infrastructure/repository/tsdb"
)

type PortfolioItemReplicaSet struct {
	*ReplicaSet
}

var _ portfolio_item.ReplicaSet = (*PortfolioItemReplicaSet)(nil)

func NewPortfolioItemReplicaSet(replicaSet *ReplicaSet) *PortfolioItemReplicaSet {
	return &PortfolioItemReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *PortfolioItemReplicaSet) WriteRepo() portfolio_item.WriteRepository {
	return tsdb.NewPortfolioItemRepository(c.ReplicaSet.WriteRepo())
}

func (c *PortfolioItemReplicaSet) ReadRepo() portfolio_item.ReadRepository {
	return tsdb.NewPortfolioItemRepository(c.ReplicaSet.ReadRepo())
}
