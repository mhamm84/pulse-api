package data

import (
	"github.com/shopspring/decimal"
	"time"
)

type Summary struct {
	Name       string                 `json:"name"`
	LastUpdate time.Time              `json:"lastUpdate"`
	Value      decimal.Decimal        `json:"value"`
	Change     decimal.Decimal        `json:"change"`
	Slug       string                 `json:"slug"`
	Extras     map[string]interface{} `json:"extras"`
}
