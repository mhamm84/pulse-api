package economic

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	serviceTimeout = 10
)

type AlphaVantageEconomicResponse struct {
	Name     string                     `json:"Name"`
	Interval string                     `json:"Interval"`
	Unit     string                     `json:"Unit"`
	Data     []AlphaVantageEconomicData `json:"Data"`
}

type AlphaVantageEconomicData struct {
	Date  time.Time `json:"date,string"`
	Value big.Float `json:"value,string"`
}

func (l *AlphaVantageEconomicData) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "date" {
			t, err := time.Parse("2006-01-02", v)
			if err != nil {
				return err
			}
			l.Date = t
		}
		if strings.ToLower(k) == "value" {
			fv, err := strconv.ParseFloat(v, 64)
			v := big.NewFloat(fv)
			if err != nil {
				return err
			}
			l.Value = *v
		}
	}
	return nil
}
