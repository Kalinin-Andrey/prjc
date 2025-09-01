package cmc_api

import (
	"info/internal/domain/concentration"
	"info/internal/domain/currency"
	"info/internal/domain/portfolio_item"
	"info/internal/domain/price_and_cap"
	"info/internal/pkg/apperror"
	"strconv"
	"time"
)

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      string `json:"elapsed"`
	CreditCount  uint   `json:"credit_count"`
}

type GetCurrencyResponse struct {
	Data   *CurrencyData `json:"data"`
	Status Status        `json:"status"`
}

type CurrencyData struct {
	ID     uint   `json:"id"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
}

func (e *CurrencyData) Currency() *currency.Currency {
	return &currency.Currency{
		ID:     e.ID,
		Symbol: e.Symbol,
		Slug:   e.Slug,
		Name:   e.Name,
	}
}

type DetailChartResponse struct {
	Data   *DetailChartData `json:"data"`
	Status Status           `json:"status"`
}

type DetailChartData struct {
	CurrencyID uint
	Points     *DetailChartPoints `json:"points"`
}

func (e *DetailChartData) PriceAndCapList() (*price_and_cap.PriceAndCapList, error) {
	if e == nil || e.Points == nil || len(*e.Points) == 0 || e.CurrencyID == 0 {
		return nil, apperror.ErrNotFound
	}
	var t string
	var point DetailChartPoint
	res := make(price_and_cap.PriceAndCapList, 0, len(*e.Points))
	for t, point = range *e.Points {
		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, err
		}

		if len(point.V) < 3 {
			continue
		}

		res = append(res, price_and_cap.PriceAndCap{
			CurrencyID:  e.CurrencyID,
			Price:       point.V[0],
			DailyVolume: point.V[1],
			Cap:         point.V[2],
			Ts:          time.Unix(i, 0),
		})
	}

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}
	return &res, nil
}

type DetailChartPoints map[string]DetailChartPoint

type DetailChartPoint struct {
	V []float64 `json:"v"`
	C []float64 `json:"c"`
}

type GetAnalyticsResponse struct {
	Data   *GetAnalyticsData `json:"data"`
	Status Status            `json:"status"`
}

type GetAnalyticsData struct {
	CurrencyID              uint
	AddressesByHoldings     *interface{}             `json:"addressesByHoldings"`
	HistoricalConcentration *HistoricalConcentration `json:"historicalConcentration"`
	AddressByTimeHeld       *interface{}             `json:"addressByTimeHeld"`
	AverageTransactionFees  *interface{}             `json:"averageTransactionFees"`
}

type HistoricalConcentration struct {
	HistoricalConcentrationDetails *[]HistoricalConcentrationDetailsPoint `json:"historicalConcentrationDetails"`
	HistoricalConcentrationAgg     *HistoricalConcentrationAgg            `json:"historicalConcentrationAgg"`
}

type HistoricalConcentrationAgg struct {
	WhalesPercent   float64 `json:"whalesPercent"`
	WhalesVolume    float64 `json:"whalesVolume"`
	OthersPercent   float64 `json:"othersPercent"`
	OthersVolume    float64 `json:"othersVolume"`
	InvestorPercent float64 `json:"investorPercent"`
	InvestorVolume  float64 `json:"investorVolume"`
	RetailPercent   float64 `json:"retailPercent"`
	RetailVolume    float64 `json:"retailVolume"`
}

type HistoricalConcentrationDetailsPoint struct {
	Date      string  `json:"date"`
	Whales    float64 `json:"whales"`
	Others    float64 `json:"others"`
	Retail    float64 `json:"retail"`
	Investors float64 `json:"investors"`
}

func (e *GetAnalyticsData) ConcentrationList() (*concentration.ConcentrationList, error) {
	if e == nil || e.HistoricalConcentration == nil || e.HistoricalConcentration.HistoricalConcentrationDetails == nil || len(*e.HistoricalConcentration.HistoricalConcentrationDetails) == 0 || e.CurrencyID == 0 {
		return &concentration.ConcentrationList{}, nil
	}
	return e.HistoricalConcentration.ConcentrationList(e.CurrencyID)
}

func (e *HistoricalConcentration) ConcentrationList(currencyID uint) (*concentration.ConcentrationList, error) {
	if e == nil || e.HistoricalConcentrationDetails == nil || len(*e.HistoricalConcentrationDetails) == 0 || currencyID == 0 {
		return nil, apperror.ErrNotFound
	}

	res := make(concentration.ConcentrationList, 0, len(*e.HistoricalConcentrationDetails))
	var item HistoricalConcentrationDetailsPoint
	var t time.Time
	var err error
	for _, item = range *e.HistoricalConcentrationDetails {
		if t, err = time.Parse(time.DateOnly, item.Date); err != nil {
			return nil, err
		}

		res = append(res, concentration.Concentration{
			CurrencyID: currencyID,
			Whales:     item.Whales,
			Investors:  item.Investors,
			Retail:     item.Retail,
			D:          t,
		})
	}
	return &res, nil
}

type GetPortfolioSummaryRequest struct {
	PortfolioSourceId string `json:"portfolioSourceId"`
	PortfolioType     string `json:"portfolioType"`
	CryptoUnit        uint   `json:"cryptoUnit"`
	CurrentPage       uint   `json:"currentPage"`
	PageSize          uint   `json:"pageSize"`
}

type GetPortfolioSummaryResponse struct {
	Data   *PortfolioSummary `json:"data"`
	Status Status            `json:"status"`
}

type PortfolioSummary struct {
	PortfolioType string             `json:"portfolioType"`
	ManualSummary []PortfolioContent `json:"manualSummary"`
}

type PortfolioContent struct {
	CurrentPage uint              `json:"currentPage"`
	List        PortfolioItemList `json:"list"`
}

type PortfolioItemList []PortfolioItem

func (l *PortfolioItemList) PortfolioItemList() *portfolio_item.PortfolioItemList {
	if l == nil || len(*l) == 0 {
		return nil
	}
	res := make(portfolio_item.PortfolioItemList, 0, len(*l))
	var item PortfolioItem
	var resItem *portfolio_item.PortfolioItem

	for _, item = range *l {
		resItem = item.PortfolioItem()
		res = append(res, *resItem)
	}

	return &res
}

func (l *PortfolioItemList) SetPortfolioSourceId(portfolioSourceId string) *PortfolioItemList {
	if l == nil || len(*l) == 0 {
		return nil
	}
	var i int
	for i = range *l {
		(*l)[i].PortfolioSourceID = portfolioSourceId
	}
	return l
}

type PortfolioItem struct {
	PortfolioSourceID string    `json:"portfolioSourceId"`
	CurrencyID        uint      `json:"cryptocurrencyId"`
	Amount            float64   `json:"amount"`
	CurrentPrice      float64   `json:"currentPrice"`
	CryptoHoldings    float64   `json:"cryptoHoldings"`
	HoldingsPercent   float64   `json:"holdingsPercent"`
	BuyAvgPrice       float64   `json:"buyAvgPrice"`
	PlPercentValue    float64   `json:"plPercentValue"`
	PlValue           float64   `json:"plValue"`
	TotalBuySpent     float64   `json:"totalBuySpent"`
	UpdatedAt         time.Time `json:"lastUpdated"`
}

func (e *PortfolioItem) PortfolioItem() *portfolio_item.PortfolioItem {
	return &portfolio_item.PortfolioItem{
		PortfolioSourceID: e.PortfolioSourceID,
		CurrencyID:        e.CurrencyID,
		Amount:            e.Amount,
		CurrentPrice:      e.CurrentPrice,
		CryptoHoldings:    e.CryptoHoldings,
		HoldingsPercent:   e.HoldingsPercent,
		BuyAvgPrice:       e.BuyAvgPrice,
		PlPercentValue:    e.PlPercentValue,
		PlValue:           e.PlValue,
		TotalBuySpent:     e.TotalBuySpent,
		UpdatedAt:         e.UpdatedAt,
	}
}
