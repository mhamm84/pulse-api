package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/mhamm84/pulse-api/internal/data/economic"
	"github.com/mhamm84/pulse-api/internal/validator"
	"net/http"
	"strings"
	"time"
)

func (app *application) cpiDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(app, economic.CPI, w, r)
}

func (app *application) consumerSentimentDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(app, economic.ConsumerSentiment, w, r)
}

func (app *application) retailSalesDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(app, economic.RetailSales, w, r)
}

func (app *application) treasuryYieldByYears(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	maturity := params.ByName("maturity")
	reportType := economic.ReportTypeTreasuryYieldMaturity(maturity)
	if reportType == economic.Unknown {
		app.logger.PrintInfo("Unknown maturity to get treasury yield data", map[string]interface{}{
			"maturity": maturity,
		})
		app.badRequestHandler(w, r)
		return
	}
	getEconomicDataByYears(app, reportType, w, r)
}

func getEconomicDataByYears(app *application, report economic.ReportType, w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	qs := r.URL.Query()
	years := app.readInt(qs, "years", time.Now().Year(), v)

	v.Check(years > 0, fmt.Sprintf("%s.years", report), "years must be a positive value")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	data, err := app.services.alphaVantageEconomicService.GetIntervalWithPercentChange(report, years)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.WriteJson(w, http.StatusOK, envelope{strings.ToLower(report.String()): data}, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}
}
