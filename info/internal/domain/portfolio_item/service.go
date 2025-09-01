package portfolio_item

import (
	"context"
	"errors"
	"fmt"
	"info/internal/pkg/apperror"
	"runtime/debug"
	"time"
)

type CmcApi interface {
	GetPortfolioSummary(ctx context.Context, portfolioSourceId string) (*PortfolioItemList, error)
}

type Service struct {
	replicaSet ReplicaSet
	cmcApi     CmcApi
}

func NewService(replicaSet ReplicaSet, cmcApi CmcApi) *Service {
	return &Service{
		replicaSet: replicaSet,
		cmcApi:     cmcApi,
	}
}

const (
	defaultCapacity = 100
)

func (s *Service) MGetByPortfolioSourceId(ctx context.Context, portfolioSourceId uint) (*PortfolioItemMap, error) {
	return s.replicaSet.ReadRepo().MGetByPortfolioSourceId(ctx, portfolioSourceId)
}

func (s *Service) mUpsert(ctx context.Context, entities *PortfolioItemList) error {
	return s.replicaSet.WriteRepo().MUpsert(ctx, entities)
}

func (s *Service) Import(ctx context.Context, portfolioSourceIDs *[]string) (err error) {
	const metricName = "portfolio_item.Service.Import"

	if portfolioSourceIDs == nil || len(*portfolioSourceIDs) == 0 {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Recover from panic: %v; stacktrace from panic: %s", apperror.ErrInternal, r, string(debug.Stack())))
		}
	}()

	var portfolioSourceID string
	for _, portfolioSourceID = range *portfolioSourceIDs {
		if err = s.importItem(ctx, portfolioSourceID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) importItem(ctx context.Context, portfolioSourceID string) (err error) {
	const metricName = "portfolio_item.Service.importItem"

	if portfolioSourceID == "" {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Recover from panic: %v; stacktrace from panic: %s", apperror.ErrInternal, r, string(debug.Stack())))
		}
	}()

	time.Sleep(3 * time.Second)
	l, err := s.cmcApi.GetPortfolioSummary(ctx, portfolioSourceID)
	if err != nil {
		return fmt.Errorf("[%w] cmcApi.GetPortfolioSummary error: %w", apperror.ErrInternal, err)
	}

	return s.replicaSet.WriteRepo().MUpsert(ctx, l)
}
