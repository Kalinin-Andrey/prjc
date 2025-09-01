package price_and_cap

import "time"

const ()

type PriceAndCap struct {
	CurrencyID  uint
	Price       float64
	DailyVolume float64
	Cap         float64
	Ts          time.Time
}

func (e *PriceAndCap) Validate() error {
	return nil
}

type PriceAndCapList []PriceAndCap

func (l *PriceAndCapList) Slice() *[]PriceAndCap {
	if l == nil {
		return nil
	}
	res := []PriceAndCap(*l)
	return &res
}

func (l *PriceAndCapList) MaxTime() *time.Time {
	if l == nil || len(*l) == 0 {
		return nil
	}
	max := (*l)[0].Ts
	var item PriceAndCap
	for _, item = range *l {
		if item.Ts.After(max) {
			max = item.Ts
		}
	}
	return &max
}

func (l *PriceAndCapList) AvgInDay(d time.Time) *PriceAndCap {
	if l == nil || len(*l) == 0 || d.IsZero() {
		return nil
	}
	d1 := d.Add(-1 * time.Second)
	d2 := d.Add(time.Hour * 24)
	dayList := make(PriceAndCapList, 0, defaultCapacity)

	var item PriceAndCap
	for _, item = range *l {
		if item.Ts.After(d1) && d2.After(item.Ts) {
			dayList = append(dayList, item)
		}
	}

	return dayList.Avg()
}

func (l *PriceAndCapList) Avg() *PriceAndCap {
	if l == nil || len(*l) == 0 {
		return nil
	}
	var ts int64
	var item PriceAndCap
	le := int64(len(*l))
	leFl := float64(le)
	res := PriceAndCap{
		CurrencyID: (*l)[0].CurrencyID,
	}

	for _, item = range *l {
		ts += item.Ts.Unix() / le
		res.Price += item.Price / leFl
		res.DailyVolume += item.DailyVolume / leFl
		res.Cap += item.Cap / leFl
	}
	res.Ts = time.Unix(ts, 0)

	return &res
}

type PriceAndCapMap map[uint]PriceAndCapList
