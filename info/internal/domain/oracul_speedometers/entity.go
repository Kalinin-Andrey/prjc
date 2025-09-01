package oracul_speedometers

import (
	"time"
)

const ()

type OraculSpeedometers struct {
	CurrencyID        uint
	WhalesBuyRate     float64
	WhalesSellRate    float64
	WhalesVolume      float64
	InvestorsBuyRate  float64
	InvestorsSellRate float64
	InvestorsVolume   float64
	RetailersBuyRate  float64
	RetailersSellRate float64
	RetailersVolume   float64
	Ts                time.Time
}

func (e *OraculSpeedometers) Validate() error {
	return nil
}
