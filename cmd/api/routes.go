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
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/cpi"), app.cpiDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/consumersentiment"), app.consumerSentimentDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/retailsales"), app.retailSalesDataByYears)
	router.HandlerFunc(http.MethodGet, pathWithVersion("/%s/economic/treasuryYield/:maturity"), app.treasuryYieldByYears)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}

func pathWithVersion(pathFmt string) string {
	return fmt.Sprintf(pathFmt, apiVersion)
}
