package oracul_holder_stats

import (
	"time"
)

const ()

type OraculHolderStats struct {
	CurrencyID            uint
	WhalesVolume          float64
	WhalesTotalHolders    uint
	InvestorsVolume       float64
	InvestorsTotalHolders uint
	RetailersVolume       float64
	RetailersTotalHolders uint
	Ts                    time.Time
}

func (e *OraculHolderStats) Validate() error {
	return nil
}
