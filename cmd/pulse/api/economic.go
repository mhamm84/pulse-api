package api

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/validator"
	"net/http"
	"time"
)

func (app *application) federalFundsRate(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.FederalFundsRate, w, r)
}

func (app *application) realGdpPerCapitaDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.RealGdpPerCapita, w, r)
}

func (app *application) realGdpDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.RealGDP, w, r)
}

func (app *application) cpiDataByYears(w http.ResponseWriter, r *http.Request) {
	getEconomicDataByYears(r.Context(), app, data.CPI, w, r)
}

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
		app.badRequestHandler(w, r)
		return
	}
	getEconomicDataByYears(r.Context(), app, reportType, w, r)
}

func getEconomicDataByYears(ctx context.Context, app *application, report data.ReportType, w http.ResponseWriter, r *http.Request) {
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

	dataChan := make(chan data.EconomicWithChangeResult)
	errChan := make(chan error)

	go app.services.alphaVantageEconomicService.GetIntervalWithPercentChange(r.Context(), dataChan, errChan, report, years, input.Paging)

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
