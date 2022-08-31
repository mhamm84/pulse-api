package api

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/validator"
	"net/http"
)

const (
	yearsParam          = "years"
	timeBucketDaysParam = "timeBucketDays"
	pageParam           = "page"
	pageSizeParam       = "pageSize"
)

func (app *application) inflationExpectation(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.InflationExpectation, w, r)
}

func (app *application) inflation(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.Inflation, w, r)
}

func (app *application) nonfarmPayroll(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.NonfarmPayroll, w, r)
}

func (app *application) unemployemnt(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.Unemployment, w, r)
}

func (app *application) durableGoodsOrders(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.DurableGoodsOrders, w, r)
}

func (app *application) federalFundsRate(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.FederalFundsRate, w, r)
}

func (app *application) realGdpPerCapitaDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.RealGdpPerCapita, w, r)
}

func (app *application) realGdpDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.RealGDP, w, r)
}

//###############################################################################################
// CPI
//###############################################################################################
func (app *application) cpiDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.CPI, w, r)
}

func (app *application) cpiStats(w http.ResponseWriter, r *http.Request) {
	getStats(r.Context(), app, data.CPI, w, r)
}

//###############################################################################################
// Consumer Sentiment
//###############################################################################################
func (app *application) consumerSentimentDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.ConsumerSentiment, w, r)
}

func (app *application) retailSalesDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.RetailSales, w, r)
}

func (app *application) treasuryYieldByYears(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	maturity := params.ByName("maturity")
	reportType := data.ReportTypeTreasuryYieldMaturity(maturity)
	if reportType == data.Unknown {
		app.logger.PrintInfo("Unknown maturity to get treasury yield data", map[string]interface{}{
			"maturity": maturity,
		})
		app.badRequestResponse(w, r)
		return
	}
	getEconomicDataByYears(r.Context(), app, reportType, w, r)
}

func getStats(ctx context.Context, app *application, report data.ReportType, w http.ResponseWriter, r *http.Request) {
	var input struct {
		Years          int
		TimeBucketDays int
		Paging         data.Paging
	}

	v := validator.New()

	qs := r.URL.Query()
	years := app.readInt(qs, yearsParam, 10, v)
	input.Years = years
	timeBucketDays := app.readInt(qs, timeBucketDaysParam, 365, v)
	input.TimeBucketDays = timeBucketDays
	page := app.readInt(qs, pageParam, 1, v)
	input.Paging.Page = page
	pageSize := app.readInt(qs, pageSizeParam, 10, v)
	input.Paging.PageSize = pageSize

	data.ValidatePaging(v, input.Paging)
	checkYears(years, report, v)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	dataChan := make(chan data.EconomicStatsResult)
	errChan := make(chan error)

	go app.services.AlphaVantageEconomicService.GetStats(r.Context(), dataChan, errChan, report, years, timeBucketDays, input.Paging)

	select {
	case data := <-dataChan:
		err := app.WriteJson(w, http.StatusOK, envelope{
			"data": data.Data,
			"meta": data.Meta,
		}, nil)
		if err != nil {
			app.logger.PrintError(err, nil)
			app.serverErrorResponse(w, r, err)
		}
		return
	case err := <-errChan:
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
		return
	case <-ctx.Done():
		app.logger.PrintError(ctx.Err(), nil)
		app.serverErrorResponse(w, r, ctx.Err())
		return
	}
}

func checkYears(years int, report data.ReportType, v *validator.Validator) {
	v.Check(years > 0, fmt.Sprintf("%s.years", report), "years must be a positive value")
}

func getEconomicDataByYears(ctx context.Context, app *application, report data.ReportType, w http.ResponseWriter, r *http.Request) {

	var input struct {
		Years  int
		Paging data.Paging
	}

	v := validator.New()

	qs := r.URL.Query()
	years := app.readInt(qs, yearsParam, 10, v)
	input.Years = years

	page := app.readInt(qs, pageParam, 1, v)
	input.Paging.Page = page
	pageSize := app.readInt(qs, pageSizeParam, 10, v)
	input.Paging.PageSize = pageSize

	data.ValidatePaging(v, input.Paging)

	v.Check(years > 0, fmt.Sprintf("%s.years", report), "years must be a positive value")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	dataChan := make(chan data.EconomicWithChangeResult)
	errChan := make(chan error)

	go app.services.AlphaVantageEconomicService.GetIntervalWithPercentChange(r.Context(), dataChan, errChan, report, years, input.Paging)

	select {
	case data := <-dataChan:
		err := app.WriteJson(w, http.StatusOK, envelope{
			"data": data.Data,
			"meta": data.Meta,
		}, nil)

		if err != nil {
			app.logger.PrintError(err, nil)
			app.serverErrorResponse(w, r, err)
		}
		return
	case err := <-errChan:
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
		return
	case <-ctx.Done():
		app.logger.PrintError(ctx.Err(), nil)
		app.serverErrorResponse(w, r, ctx.Err())
		return
	}
}
