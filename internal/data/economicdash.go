package data

import (
	"github.com/shopspring/decimal"
	"time"
)

type SummaryHeader struct {
	HeaderName string    `json:"headerName"`
	Summaries  []Summary `json:"summaries"`
}

type Summary struct {
	Name       string          `json:"name"`
	LastUpdate time.Time       `json:"lastUpdate"`
	Value      decimal.Decimal `json:"value"`
	Change     decimal.Decimal `json:"change"`
}
