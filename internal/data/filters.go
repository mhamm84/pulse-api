package data

import (
	"github.com/mhamm84/pulse-api/internal/validator"
	"math"
)

type Paging struct {
	Page     int
	PageSize int
}

func (p *Paging) Limit() int {
	return p.PageSize
}

func (p *Paging) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func ValidatePaging(v *validator.Validator, p Paging) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "must be a maximum of 100")
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {

		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
