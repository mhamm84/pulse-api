package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const apiVersion = "v1"

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/healthcheck"), app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/dashboard"), app.economicDashHandler)

	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/real_gdp"), app.realGdpDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/real_gdp_per_capita"), app.realGdpPerCapitaDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/cpi"), app.cpiDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/consumer_sentiment"), app.consumerSentimentDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/retail_sales"), app.retailSalesDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/treasury_yield/:maturity"), app.treasuryYieldByYears)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}

func pathWithVersion(pathFmt string) string {
	return fmt.Sprintf(pathFmt, apiVersion)
}
