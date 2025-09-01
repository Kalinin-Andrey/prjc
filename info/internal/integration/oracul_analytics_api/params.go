package oracul_analytics_api

import (
	"info/internal/domain/oracul_analytics"
	"info/internal/domain/oracul_daily_balance_stats"
	"info/internal/domain/oracul_holder_stats"
	"info/internal/domain/oracul_speedometers"
	"strconv"
	"time"
)

type GetHoldersStatsResponse struct {
	WhalesConcentration string                `json:"whales_concentration"`
	WormIndex           string                `json:"worm_index"`
	GrowthFuel          string                `json:"growth_fuel"`
	Speedometers        *Speedometers         `json:"speedometers"`
	HolderStats         *HolderStats          `json:"holder_stats"`
	DailyBalanceStats   DailyBalanceStatsList `json:"daily_balance_stats"`
}

func (r *GetHoldersStatsResponse) ImportData(currencyID uint, ts time.Time) (res *oracul_analytics.ImportData, err error) {

	res = &oracul_analytics.ImportData{}

	if res.OraculAnalytics, err = r.OraculAnalytics(currencyID, ts); err != nil {
		return nil, err
	}
	if res.OraculSpeedometers, err = r.Speedometers.OraculSpeedometers(currencyID, ts); err != nil {
		return nil, err
	}
	if res.OraculHolderStats, err = r.HolderStats.OraculHolderStats(currencyID, ts); err != nil {
		return nil, err
	}
	if res.OraculDailyBalanceStatsList, err = r.DailyBalanceStats.OraculDailyBalanceStatsList(currencyID); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *GetHoldersStatsResponse) OraculAnalytics(currencyID uint, ts time.Time) (res *oracul_analytics.OraculAnalytics, err error) {
	res = &oracul_analytics.OraculAnalytics{
		CurrencyID: currencyID,
		Ts:         ts,
	}

	if res.WhalesConcentration, err = strconv.ParseFloat(r.WhalesConcentration, 64); err != nil {
		return nil, err
	}
	if res.WormIndex, err = strconv.ParseFloat(r.WormIndex, 64); err != nil {
		return nil, err
	}
	if res.GrowthFuel, err = strconv.ParseFloat(r.GrowthFuel, 64); err != nil {
		return nil, err
	}

	return res, nil
}

type Speedometers struct {
	Whales    *SpeedometersItem `json:"whales"`
	Investors *SpeedometersItem `json:"investors"`
	Retailers *SpeedometersItem `json:"retailers"`
}
type SpeedometersItem struct {
	BuyRate  string `json:"buy_rate"`
	SellRate string `json:"sell_rate"`
	Volume   string `json:"volume"`
}

func (e *Speedometers) OraculSpeedometers(currencyID uint, ts time.Time) (res *oracul_speedometers.OraculSpeedometers, err error) {
	res = &oracul_speedometers.OraculSpeedometers{
		CurrencyID: currencyID,
		Ts:         ts,
	}

	if res.WhalesBuyRate, err = strconv.ParseFloat(e.Whales.BuyRate, 64); err != nil {
		return nil, err
	}
	if res.WhalesSellRate, err = strconv.ParseFloat(e.Whales.SellRate, 64); err != nil {
		return nil, err
	}
	if res.WhalesVolume, err = strconv.ParseFloat(e.Whales.Volume, 64); err != nil {
		return nil, err
	}
	if res.InvestorsBuyRate, err = strconv.ParseFloat(e.Investors.BuyRate, 64); err != nil {
		return nil, err
	}
	if res.InvestorsSellRate, err = strconv.ParseFloat(e.Investors.SellRate, 64); err != nil {
		return nil, err
	}
	if res.InvestorsVolume, err = strconv.ParseFloat(e.Investors.Volume, 64); err != nil {
		return nil, err
	}
	if res.RetailersBuyRate, err = strconv.ParseFloat(e.Retailers.BuyRate, 64); err != nil {
		return nil, err
	}
	if res.RetailersSellRate, err = strconv.ParseFloat(e.Retailers.SellRate, 64); err != nil {
		return nil, err
	}
	if res.RetailersVolume, err = strconv.ParseFloat(e.Retailers.Volume, 64); err != nil {
		return nil, err
	}

	return res, nil
}

type HolderStats struct {
	Whales    *HolderStatsItem `json:"whales"`
	Investors *HolderStatsItem `json:"investors"`
	Retailers *HolderStatsItem `json:"retailers"`
}
type HolderStatsItem struct {
	Volume       string `json:"volume"`
	TotalHolders uint   `json:"total_holders"`
}

func (e *HolderStats) OraculHolderStats(currencyID uint, ts time.Time) (res *oracul_holder_stats.OraculHolderStats, err error) {
	res = &oracul_holder_stats.OraculHolderStats{
		CurrencyID:            currencyID,
		WhalesTotalHolders:    e.Whales.TotalHolders,
		InvestorsTotalHolders: e.Investors.TotalHolders,
		RetailersTotalHolders: e.Retailers.TotalHolders,
		Ts:                    ts,
	}

	if res.WhalesVolume, err = strconv.ParseFloat(e.Whales.Volume, 64); err != nil {
		return nil, err
	}
	if res.InvestorsVolume, err = strconv.ParseFloat(e.Investors.Volume, 64); err != nil {
		return nil, err
	}
	if res.RetailersVolume, err = strconv.ParseFloat(e.Retailers.Volume, 64); err != nil {
		return nil, err
	}

	return res, nil
}

type DailyBalanceStatsList map[string]DailyBalanceStats
type DailyBalanceStats struct {
	Whales    *DailyBalanceStatsItem `json:"whales"`
	Investors *DailyBalanceStatsItem `json:"investors"`
	Retailers *DailyBalanceStatsItem `json:"retailers"`
}
type DailyBalanceStatsItem struct {
	Balance      string `json:"balance"`
	TotalHolders uint   `json:"total_holders"`
}

func (e *DailyBalanceStatsList) OraculDailyBalanceStatsList(currencyID uint) (*oracul_daily_balance_stats.OraculDailyBalanceStatsList, error) {
	if e == nil || len(*e) == 0 {
		return nil, nil
	}
	var err error
	var date string
	var d time.Time
	var item DailyBalanceStats
	var resItem *oracul_daily_balance_stats.OraculDailyBalanceStats
	res := make(oracul_daily_balance_stats.OraculDailyBalanceStatsList, 0, len(*e))

	for date, item = range *e {
		d, err = time.Parse(time.DateOnly, date)
		if err != nil {
			return nil, err
		}
		if resItem, err = item.OraculDailyBalanceStats(currencyID, d); err != nil {
			return nil, err
		}
		res = append(res, *resItem)
	}

	return &res, nil
}

func (e *DailyBalanceStats) OraculDailyBalanceStats(currencyID uint, d time.Time) (res *oracul_daily_balance_stats.OraculDailyBalanceStats, err error) {
	res = &oracul_daily_balance_stats.OraculDailyBalanceStats{
		CurrencyID:            currencyID,
		WhalesTotalHolders:    e.Whales.TotalHolders,
		InvestorsTotalHolders: e.Investors.TotalHolders,
		RetailersTotalHolders: e.Retailers.TotalHolders,
		D:                     d,
	}

	if res.WhalesBalance, err = strconv.ParseFloat(e.Whales.Balance, 64); err != nil {
		return nil, err
	}
	if res.InvestorsBalance, err = strconv.ParseFloat(e.Investors.Balance, 64); err != nil {
		return nil, err
	}
	if res.RetailersBalance, err = strconv.ParseFloat(e.Retailers.Balance, 64); err != nil {
		return nil, err
	}

	return res, nil
}
