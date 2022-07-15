package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/validator"
	"net/http"
	"time"
)

func (app *application) cpiDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(app, data.CPI, w, r)
}

func (app *application) consumerSentimentDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(app, data.ConsumerSentiment, w, r)
}

func (app *application) retailSalesDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(app, data.RetailSales, w, r)
}

func (app *application) treasuryYieldByYears(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	maturity := params.ByName("maturity")
	reportType := data.ReportTypeTreasuryYieldMaturity(maturity)
	if reportType == data.Unknown {
		app.logger.PrintInfo("Unknown maturity to get treasury yield data", map[string]interface{}{
			"maturity": maturity,
		})
		app.badRequestHandler(w, r)
		return
	}
	getEconomicDataByYears(app, reportType, w, r)
}

func getEconomicDataByYears(app *application, report data.ReportType, w http.ResponseWriter, r *http.Request) {
	var input struct {
		Years  int
		Paging data.Paging
	}

	v := validator.New()

	qs := r.URL.Query()
	years := app.readInt(qs, "years", time.Now().Year(), v)
	input.Years = years

	page := app.readInt(qs, "page", 1, v)
	input.Paging.Page = page
	pageSize := app.readInt(qs, "page_size", 10, v)
	input.Paging.PageSize = pageSize

	data.ValidatePaging(v, input.Paging)

	v.Check(years > 0, fmt.Sprintf("%s.years", report), "years must be a positive value")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	data, meta, err := app.services.alphaVantageEconomicService.GetIntervalWithPercentChange(report, years, input.Paging)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.WriteJson(w, http.StatusOK, envelope{
		"data": data,
		"meta": meta,
	}, nil)

	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}
}
