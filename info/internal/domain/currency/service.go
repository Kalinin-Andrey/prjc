package currency

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"info/internal/domain"
	"info/internal/domain/concentration"
	"info/internal/domain/oracul_analytics"
	"info/internal/domain/price_and_cap"
	"info/internal/pkg/apperror"
	"math"
	"runtime/debug"
	"time"
)

const (
	defaultCapacity = 100
)

type CmcApi interface {
	GetCurrency(ctx context.Context, currencySlug string) (*Currency, error)
}

type CmcProApi interface {
	GetCurrenciesBySlugs(ctx context.Context, slugs *[]string) (currencyMap CurrencyMap, err error)
}

type Service struct {
	replicaSet      ReplicaSet
	priceAndCap     *price_and_cap.Service
	concentration   *concentration.Service
	oraculAnalytics *oracul_analytics.Service
	cmcApi          CmcApi
	cmcProApi       CmcProApi
}

func NewService(replicaSet ReplicaSet, priceAndCap *price_and_cap.Service, concentration *concentration.Service, oraculAnalytics *oracul_analytics.Service, cmcApi CmcApi, cmcProApi CmcProApi) *Service {
	return &Service{
		replicaSet:      replicaSet,
		priceAndCap:     priceAndCap,
		concentration:   concentration,
		oraculAnalytics: oraculAnalytics,
		cmcApi:          cmcApi,
		cmcProApi:       cmcProApi,
	}
}

func (s *Service) Create(ctx context.Context, entity *Currency) (ID uint, err error) {
	return s.replicaSet.WriteRepo().Create(ctx, entity)
}

func (s *Service) Update(ctx context.Context, entity *Currency) error {
	return s.replicaSet.WriteRepo().Update(ctx, entity)
}

func (s *Service) Delete(ctx context.Context, ID uint) error {
	return s.replicaSet.WriteRepo().Delete(ctx, ID)
}

func (s *Service) Get(ctx context.Context, ID uint) (*Currency, error) {
	return s.replicaSet.ReadRepo().Get(ctx, ID)
}

func (s *Service) GetAll(ctx context.Context) (*CurrencyList, error) {
	return s.replicaSet.ReadRepo().GetAll(ctx)
}

func (s *Service) Import(ctx context.Context, listOfCurrencySlugs *[]string) (err error) {
	const metricName = "currency.Service.ImportTx"
	var tx domain.Tx

	currencyList, err := s.baseImport(ctx, listOfCurrencySlugs)
	if err != nil {
		return err
	}

	tx, err = s.replicaSet.WriteRepo().Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Recover from panic: %v; stacktrace from panic: %s", apperror.ErrInternal, r, string(debug.Stack())))
		}

		if err == nil {
			if err = tx.Commit(ctx); err == nil {
				return
			}
			err = fmt.Errorf("[%w] "+metricName+" Commit error: %w", apperror.ErrInternal, err)
		}

		if tx != nil {
			if err2 := tx.Rollback(ctx); err2 != nil {
				err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Rollback error: %w", apperror.ErrInternal, err))
			}
		}

	}()

	importMaxTimeMap, err := s.replicaSet.WriteRepo().GetImportMaxTimeForUpdateTx(ctx, tx, currencyList.IDs())
	if err != nil {
		return err
	}

	var importMaxTimeItem ImportMaxTime
	var currency Currency
	var ok bool
	var i int
	for i, currency = range *currencyList {
		importMaxTimeItem, ok = importMaxTimeMap[currency.ID]
		if !ok {
			importMaxTimeItem = ImportMaxTime{
				CurrencyID: currency.ID,
			}
		}
		fmt.Printf("%d. %s\n", i, currency.Symbol)
		time.Sleep(6 * time.Second)
		if importMaxTimeItem.PriceAndCap, err = s.priceAndCap.ImportTx(ctx, tx, importMaxTimeItem.CurrencyID, importMaxTimeItem.PriceAndCap); err != nil {
			return err
		}
		time.Sleep(12 * time.Second)
		if importMaxTimeItem.Concentration, err = s.concentration.ImportTx(ctx, tx, importMaxTimeItem.CurrencyID, importMaxTimeItem.Concentration); err != nil {
			return err
		}

		importMaxTimeMap[currency.ID] = importMaxTimeItem
	}

	if err = s.replicaSet.WriteRepo().MUpsertImportMaxTimeMapTx(ctx, tx, importMaxTimeMap); err != nil {
		return err
	}

	//tokenAddressList, err := s.replicaSet.ReadRepo().MGetTokenAddress(ctx, currencyList.IDs())
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println("OraculAnalytics.Import")
	//if err = s.oraculAnalytics.Import(ctx, TokenAddressList2OraculAnalyticsTokenAddressList(tokenAddressList)); err != nil {
	//	return err
	//}

	return nil
}

func (s *Service) baseImport(ctx context.Context, listOfCurrencySlugs *[]string) (currencyList *CurrencyList, err error) {
	const metricName = "currency.Service.baseImport"
	if listOfCurrencySlugs == nil || len(*listOfCurrencySlugs) == 0 {
		return nil, nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Recover from panic: %v; stacktrace from panic: %s", apperror.ErrInternal, r, string(debug.Stack())))
		}
	}()

	currencyMap, err := s.cmcProApi.GetCurrenciesBySlugs(ctx, listOfCurrencySlugs)
	if err != nil {
		return nil, err
	}
	l := currencyMap.List()

	return l, s.replicaSet.WriteRepo().MUpsert(ctx, l)
}

func (s *Service) baseSimpleImport(ctx context.Context, listOfCurrencySlugs *[]string) (currencyList *CurrencyList, err error) {
	const metricName = "currency.Service.baseSimpleImport"
	if listOfCurrencySlugs == nil || len(*listOfCurrencySlugs) == 0 {
		return nil, nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("[%w] "+metricName+" Recover from panic: %v; stacktrace from panic: %s", apperror.ErrInternal, r, string(debug.Stack())))
		}
	}()

	exists, err := s.replicaSet.ReadRepo().MGetBySlug(ctx, listOfCurrencySlugs)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		l := make(CurrencyList, 0, len(*listOfCurrencySlugs))
		exists = &l
	}
	notExistsSlugs := make([]string, 0, len(*listOfCurrencySlugs)-len(*exists))
	existsSlugsMap := make(map[string]struct{})
	var item Currency
	for _, item = range *exists {
		existsSlugsMap[item.Slug] = struct{}{}
	}

	var ok bool
	var slug string
	for _, slug = range *listOfCurrencySlugs {
		if _, ok = existsSlugsMap[slug]; !ok {
			notExistsSlugs = append(notExistsSlugs, slug)
		}
	}

	if len(notExistsSlugs) > 0 {
		importedCurrencyList, err := s.baseSimpleImportBySlug(ctx, &notExistsSlugs)
		if err != nil {
			return nil, err
		}
		(*exists) = append(*exists, (*importedCurrencyList)...)
	}
	return exists, nil
}

func (s *Service) baseSimpleImportBySlug(ctx context.Context, listOfCurrencySlugs *[]string) (importedCurrencyList *CurrencyList, err error) {
	if listOfCurrencySlugs == nil || len(*listOfCurrencySlugs) == 0 {
		return nil, nil
	}

	var slug string
	var item *Currency
	importedCurrency := make(CurrencyList, 0, len(*listOfCurrencySlugs))
	for _, slug = range *listOfCurrencySlugs {
		if item, err = s.cmcApi.GetCurrency(ctx, slug); err != nil {
			return nil, err
		}
		item.IsForObserving = true
		_, err = s.replicaSet.WriteRepo().Create(ctx, item)
		if err != nil {
			return nil, err
		}
		importedCurrency = append(importedCurrency, *item)
	}
	return &importedCurrency, s.createEmptyImportMaxTime(ctx, importedCurrency.IDs())
}

func (s *Service) createEmptyImportMaxTime(ctx context.Context, IDs *[]uint) error {
	if IDs == nil || len(*IDs) == 0 {
		return nil
	}

	maxTimeList := make([]ImportMaxTime, 0, len(*IDs))
	var ID uint
	for _, ID = range *IDs {
		maxTimeList = append(maxTimeList, ImportMaxTime{
			CurrencyID: ID,
		})
	}

	return s.replicaSet.WriteRepo().MCreateImportMaxTime(ctx, &maxTimeList)
}

func (s *Service) Report_BiggestFall(ctx context.Context, limit uint) (*WhaleFallList, error) {
	l, err := s.getWhaleFallList(ctx)
	if err != nil {
		return nil, err
	}
	return l.SortByFallValueDesc().Limit(limit), nil
}

func (s *Service) Report_LongestFall(ctx context.Context, limit uint) (*WhaleFallList, error) {
	l, err := s.getWhaleFallList(ctx)
	if err != nil {
		return nil, err
	}
	return l.SortByFallDurationDesc().Limit(limit), nil
}

func (s *Service) getWhaleFallList(ctx context.Context) (*WhaleFallList, error) {
	currencyList, err := s.replicaSet.ReadRepo().GetAll(ctx)
	if err != nil {
		return nil, err
	}
	currencyIDs := currencyList.IDs()

	priceAndCapMap, err := s.priceAndCap.MGet(ctx, currencyIDs)
	if err != nil {
		return nil, err
	}

	concentrationMap, err := s.concentration.MGet(ctx, currencyIDs)
	if err != nil {
		return nil, err
	}

	return s.calcWhaleFallList(currencyList, priceAndCapMap, concentrationMap), nil
}

func (s *Service) calcWhaleFallList(currencyList *CurrencyList, priceAndCapMap price_and_cap.PriceAndCapMap, concentrationMap concentration.ConcentrationMap) *WhaleFallList {
	if currencyList == nil || priceAndCapMap == nil || concentrationMap == nil {
		return nil
	}
	var ok bool
	var i int
	var currency Currency
	var priceAndCapList price_and_cap.PriceAndCapList
	var concentrationList concentration.ConcentrationList
	var item *WhaleFall
	res := make(WhaleFallList, 0, len(*currencyList))

	for i, currency = range *currencyList {
		if i == 91 {
			fmt.Println(".")
		}
		if priceAndCapList, ok = priceAndCapMap[currency.ID]; !ok {
			continue
		}
		if concentrationList, ok = concentrationMap[currency.ID]; !ok {
			continue
		}
		item = s.calcWhaleFall(&currency, &priceAndCapList, &concentrationList)
		if item == nil {
			currencyJson, _ := json.Marshal(currency)
			fmt.Printf("calcWhaleFall: empty result for currency: %s", string(currencyJson))
			continue
		}
		res = append(res, *item)
	}

	return &res
}

func (s *Service) calcWhaleFall(currency *Currency, priceAndCapList *price_and_cap.PriceAndCapList, concentrationList *concentration.ConcentrationList) *WhaleFall {
	if currency == nil || priceAndCapList == nil || concentrationList == nil {
		return nil
	}
	const (
		maxPeriod = time.Hour * 24 * 61 // Максимальный период времени, который смотрим
		maxBreak  = time.Hour * 24 * 5  // Максимальный перерыв в тренде
	)
	now := time.Now()
	// minTime: ограничение, дальше которого не смотрим
	minTime := now.Add(-maxPeriod)
	var inFall bool
	var i int
	var prev, next, valueFrom, valueTo, localStart *concentration.Concentration
	// т.к. concentrationList отсортирован в порядке убывания по времени, в цикле next идёт перед prev
	for i = range *concentrationList {
		prev = &(*concentrationList)[i]
		// первую итерацию просто пропустим
		if i == 0 {
			next = prev
			continue
		}
		// дальше minTime не смотрим, останавливаем цикл
		if prev.D.Before(minTime) {
			break
		}

		if !inFall {
			valueTo = next
		}

		// если это не спад, то пропустим
		if valueTo.Whales >= prev.Whales {
			// если уже нашли падение, то проверяем на maxBreak
			if inFall && localStart != nil && localStart.D.Sub(prev.D) >= maxBreak {
				break
			}
			next = prev
			continue
		}
		inFall = true

		localStart = prev
		next = prev
	}
	if localStart != nil {
		valueFrom = localStart
	}
	if !inFall || valueFrom == nil || valueTo == nil {
		return nil
	}

	priceAndCapFrom := priceAndCapList.AvgInDay(valueFrom.D)
	priceAndCapTo := priceAndCapList.AvgInDay(valueTo.D)

	return &WhaleFall{
		Symbol:           currency.Symbol,
		FallDuration:     valueTo.D.Sub(valueFrom.D),
		DayFrom:          valueFrom.D,
		DayTo:            valueTo.D,
		FallValue:        valueFrom.Whales - valueTo.Whales,
		ValueFrom:        valueFrom.Whales,
		ValueTo:          valueTo.Whales,
		FallValuePercent: round(((valueTo.Whales * 100) / valueFrom.Whales)),
		FallCap:          priceAndCapFrom.Cap - priceAndCapTo.Cap,
		CapFrom:          priceAndCapFrom.Cap,
		CapTo:            priceAndCapTo.Cap,
		FallCapPercent:   round(((priceAndCapTo.Cap * 100) / priceAndCapFrom.Cap)),
		FallPrice:        priceAndCapFrom.Price - priceAndCapTo.Price,
		PriceFrom:        priceAndCapFrom.Price,
		PriceTo:          priceAndCapTo.Price,
		FallPricePercent: round(((priceAndCapTo.Price * 100) / priceAndCapFrom.Price)),
	}
}

func round(v float64) float64 {
	return math.Round(v*100) / 100
}

func TokenAddress2OraculAnalyticsTokenAddress(e *TokenAddress) *oracul_analytics.TokenAddress {
	return &oracul_analytics.TokenAddress{
		CurrencyID: e.CurrencyID,
		Blockchain: e.Blockchain,
		Address:    e.Address,
	}
}

func TokenAddressList2OraculAnalyticsTokenAddressList(l *TokenAddressList) *oracul_analytics.TokenAddressList {
	if l == nil || len(*l) == 0 {
		return nil
	}

	res := make(oracul_analytics.TokenAddressList, 0, len(*l))
	var item TokenAddress

	for _, item = range *l {
		res = append(res, *TokenAddress2OraculAnalyticsTokenAddress(&item))
	}

	return &res
}
