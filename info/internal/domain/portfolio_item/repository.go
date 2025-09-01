package portfolio_item

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	MUpsert(ctx context.Context, entities *PortfolioItemList) error
}

type ReadRepository interface {
	MGetByPortfolioSourceId(ctx context.Context, portfolioSourceId uint) (*PortfolioItemMap, error)
}
