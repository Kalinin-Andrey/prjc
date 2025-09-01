package currency

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"info/internal/pkg/apperror"
	"time"
)

const ()

type ImportMaxTime struct {
	CurrencyID    uint
	PriceAndCap   *time.Time
	Concentration *time.Time
}

type Currency struct {
	ID                            uint
	Symbol                        string
	Slug                          string
	Name                          string
	IsForObserving                bool
	CirculatingSupply             float64
	SelfReportedCirculatingSupply float64
	TotalSupply                   float64
	MaxSupply                     *float64
	LatestPrice                   float64
	CmcRank                       uint
	AddedAt                       time.Time
	Platform                      *CurrencyPlatform
	TokenAddress                  *TokenAddress
}

func (e *Currency) Validate() error {
	return nil
}

type CurrencyList []Currency

func (l *CurrencyList) IDs() *[]uint {
	if l == nil {
		return nil
	}
	res := make([]uint, 0, len(*l))
	var item Currency
	for _, item = range *l {
		res = append(res, item.ID)
	}
	return &res
}

type CurrencyMap map[uint]Currency

func (m CurrencyMap) List() *CurrencyList {
	if m == nil || len(m) == 0 {
		return nil
	}
	l := make(CurrencyList, 0, len(m))
	var item Currency

	for _, item = range m {
		l = append(l, item)
	}

	return &l
}

type CurrencyPlatform struct {
	ID           uint
	Symbol       string
	Slug         string
	Name         string
	TokenAddress string
}

func (e CurrencyPlatform) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// implement Scanner for the element type of the slice
func (e *CurrencyPlatform) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		var ok bool
		data, ok = src.([]byte)
		if !ok {
			return fmt.Errorf("[%w] type assertion to []byte failed for value: %v", apperror.ErrData, src)
		}
	}
	return json.Unmarshal(data, e)
}

type TokenAddress struct {
	CurrencyID uint
	Blockchain string
	Address    string
}

type TokenAddressList []TokenAddress
