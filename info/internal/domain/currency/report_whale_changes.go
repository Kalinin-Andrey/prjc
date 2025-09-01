package currency

import (
	"slices"
	"time"
)

type WhaleFall struct {
	Symbol           string
	FallDuration     time.Duration
	DayFrom          time.Time
	DayTo            time.Time
	FallValue        float64
	ValueFrom        float64
	ValueTo          float64
	FallValuePercent float64
	FallCap          float64
	CapFrom          float64
	CapTo            float64
	FallCapPercent   float64
	FallPrice        float64
	PriceFrom        float64
	PriceTo          float64
	FallPricePercent float64
}

type WhaleFallList []WhaleFall

func (l *WhaleFallList) SortByFallValueDesc() *WhaleFallList {
	if l == nil {
		return nil
	}
	slices.SortFunc(*l, func(a, b WhaleFall) int {
		switch {
		case a.FallValue < b.FallValue:
			return 1
		case a.FallValue > b.FallValue:
			return -1
		default:
			return 0
		}
	})
	return l
}

func (l *WhaleFallList) SortByFallDurationDesc() *WhaleFallList {
	if l == nil {
		return nil
	}
	slices.SortFunc(*l, func(a, b WhaleFall) int {
		switch {
		case a.FallDuration < b.FallDuration:
			return 1
		case a.FallDuration > b.FallDuration:
			return -1
		default:
			return 0
		}
	})
	return l
}

func (l *WhaleFallList) Limit(limit uint) *WhaleFallList {
	if l == nil {
		return nil
	}
	if uint(len(*l)) < limit {
		limit = uint(len(*l))
	}
	newList := (*l)[:limit]
	return &newList
}
