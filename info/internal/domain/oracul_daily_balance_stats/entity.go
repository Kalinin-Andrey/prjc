package oracul_daily_balance_stats

import (
	"time"
)

const ()

type OraculDailyBalanceStats struct {
	CurrencyID            uint
	WhalesBalance         float64
	WhalesTotalHolders    uint
	InvestorsBalance      float64
	InvestorsTotalHolders uint
	RetailersBalance      float64
	RetailersTotalHolders uint
	D                     time.Time
}

func (e *OraculDailyBalanceStats) Validate() error {
	return nil
}

type OraculDailyBalanceStatsList []OraculDailyBalanceStats
