package portfolio_item

import "time"

const ()

type PortfolioItem struct {
	PortfolioSourceID string
	CurrencyID        uint
	Amount            float64 // кол-во монет
	CurrentPrice      float64 // тек. цена
	CryptoHoldings    float64 // тек. общая стоимость
	HoldingsPercent   float64 // доля в портфеле
	BuyAvgPrice       float64 // ср. цена покупки
	PlPercentValue    float64 // изменение в процентах
	PlValue           float64 // изменения в $
	TotalBuySpent     float64 // всего вложено
	UpdatedAt         time.Time
}

func (e *PortfolioItem) Validate() error {
	return nil
}

type PortfolioItemList []PortfolioItem

func (l *PortfolioItemList) Slice() *[]PortfolioItem {
	if l == nil || len(*l) == 0 {
		return nil
	}
	res := []PortfolioItem(*l)
	return &res
}

func (l *PortfolioItemList) PortfoliosItemMap() PortfoliosItemMap {
	if l == nil || len(*l) == 0 {
		return nil
	}
	res := make(PortfoliosItemMap, 10)
	var item PortfolioItem
	var ok bool

	for _, item = range *l {
		if _, ok = res[item.PortfolioSourceID]; !ok {
			res[item.PortfolioSourceID] = make(PortfolioItemMap, len(*l))
		}
		res[item.PortfolioSourceID][item.CurrencyID] = item
	}

	return res
}

type PortfoliosItemMap map[string]PortfolioItemMap

func (m PortfoliosItemMap) PortfolioItemMap(portfolioSourceId string) PortfolioItemMap {
	if m == nil || len(m) == 0 {
		return nil
	}
	res, ok := m[portfolioSourceId]
	if !ok {
		return nil
	}
	return res
}

type PortfolioItemMap map[uint]PortfolioItem

func (m PortfolioItemMap) List() *PortfolioItemList {
	if m == nil || len(m) == 0 {
		return nil
	}
	res := make(PortfolioItemList, 0, len(m))
	var item PortfolioItem
	for _, item = range m {
		res = append(res, item)
	}
	return &res
}
