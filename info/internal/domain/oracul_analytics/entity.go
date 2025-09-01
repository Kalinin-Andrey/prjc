package oracul_analytics

import (
	"time"
)

const ()

type OraculAnalytics struct {
	CurrencyID          uint
	WhalesConcentration float64
	WormIndex           float64
	GrowthFuel          float64
	Ts                  time.Time
}

func (e *OraculAnalytics) Validate() error {
	return nil
}

type TokenAddress struct {
	CurrencyID uint
	Blockchain string
	Address    string
}

type TokenAddressList []TokenAddress
