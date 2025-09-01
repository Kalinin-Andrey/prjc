package concentration

import (
	"context"
	"errors"
	"fmt"
	"info/internal/domain"
	"info/internal/pkg/apperror"
	"runtime/debug"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CmcApi interface {
	GetAnalytics(ctx context.Context, CurrencyID uint, Range string) (*ConcentrationList, error)
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

	TimeRange_1M  = "1M"
	TimeRange_1Y  = "1Y"
	TimeRange_All = "All"
)

var TimeRangeList = []interface{}{
	TimeRange_1M,
	TimeRange_1Y,
	TimeRange_All,
}

func TimeRangeValidate(s string) error {
	return validation.Validate(s, validation.Required, validation.In(TimeRangeList...))
}

func (s *Service) MGet(ctx context.Context, currencyIDs *[]uint) (ConcentrationMap, error) {
	return s.replicaSet.ReadRepo().MGet(ctx, currencyIDs)
}

func (s *Service) Upsert(ctx context.Context, entity *Concentration) error {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}

func (s *Service) ImportTx(ctx context.Context, tx domain.Tx, currencyID uint, importLastTime *time.Time) (maxTime *time.Time, err error) {
	const metricName = "concentration.Service.ImportTx"
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Recover from panic: %v; stacktrace from panic: %s", apperror.ErrInternal, r, string(debug.Stack())))
		}
	}()

	if importLastTime == nil || time.Now().Add(-time.Hour*24*365).After(*importLastTime) {
		maxT, err := s.importTx(ctx, tx, currencyID, TimeRange_All)
		if err != nil {
			return nil, err
		}
		if maxTime == nil || (maxT != nil && maxT.After(*maxTime)) {
			maxTime = maxT
		}
	}

	if importLastTime == nil || time.Now().Add(-time.Hour*24*31).After(*importLastTime) {
		maxT, err := s.importTx(ctx, tx, currencyID, TimeRange_1Y)
		if err != nil {
			return nil, err
		}
		if maxTime == nil || (maxT != nil && maxT.After(*maxTime)) {
			maxTime = maxT
		}
	}

	maxT, err := s.importTx(ctx, tx, currencyID, TimeRange_1M)
	if err != nil {
		return nil, err
	}
	if maxTime == nil || (maxT != nil && maxT.After(*maxTime)) {
		maxTime = maxT
	}

	return maxTime, nil
}

func (s *Service) importTx(ctx context.Context, tx domain.Tx, currencyID uint, timeRange string) (maxTime *time.Time, err error) {
	const metricName = "concentration.Service.importTx"
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Recover from panic: %v; stacktrace from panic: %s", apperror.ErrInternal, r, string(debug.Stack())))
		}
	}()

	if err = TimeRangeValidate(timeRange); err != nil {
		return nil, err
	}

	item, err := s.cmcApi.GetAnalytics(ctx, currencyID, timeRange)
	if err != nil {
		return nil, err
	}

	if err = s.replicaSet.WriteRepo().MUpsertTx(ctx, tx, item.Slice()); err != nil {
		return nil, err
	}
	return item.MaxTime(), nil
}
