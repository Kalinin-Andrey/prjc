package concentration

import "time"

const ()

type Concentration struct {
	CurrencyID uint
	Whales     float64
	Investors  float64
	Retail     float64
	D          time.Time
}

func (e *Concentration) Validate() error {
	return nil
}

type ConcentrationList []Concentration

func (l *ConcentrationList) Slice() *[]Concentration {
	if l == nil {
		return nil
	}
	res := []Concentration(*l)
	return &res
}

func (l *ConcentrationList) MaxTime() *time.Time {
	if l == nil || len(*l) == 0 {
		return nil
	}
	max := (*l)[0].D
	var item Concentration
	for _, item = range *l {
		if item.D.After(max) {
			max = item.D
		}
	}
	return &max
}

type ConcentrationMap map[uint]ConcentrationList
