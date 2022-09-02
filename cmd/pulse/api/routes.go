package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	apiVersion         = "v1"
	economicPermission = "economic:all"
)

func (app application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/healthcheck"), app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/dashboard"), app.economicDashHandler)

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/cpi"), app.requirePermissions(economicPermission, app.cpiDataByYears))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/cpi/stats"), app.requirePermissions(economicPermission, app.cpiStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/inflation_expectation"), app.requirePermissions(economicPermission, app.inflationExpectation))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/inflation_expectation/stats"), app.requirePermissions(economicPermission, app.cpiStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/inflation"), app.requirePermissions(economicPermission, app.inflation))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/inflation/stats"), app.requirePermissions(economicPermission, app.cpiStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/nonfarm_payroll"), app.requirePermissions(economicPermission, app.nonfarmPayroll))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/nonfarm_payroll/stats"), app.requirePermissions(economicPermission, app.cpiStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/unemployment"), app.requirePermissions(economicPermission, app.unemployemnt))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/unemployment/stats"), app.requirePermissions(economicPermission, app.cpiStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/durable_goods_orders"), app.requirePermissions(economicPermission, app.durableGoodsOrders))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/durable_goods_orders/stats"), app.requirePermissions(economicPermission, app.durableGoodsOrdersStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/federal_funds_rate"), app.requirePermissions(economicPermission, app.federalFundsRate))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/federal_funds_rate/stats"), app.requirePermissions(economicPermission, app.federalFundsRateStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/real_gdp"), app.requirePermissions(economicPermission, app.realGdpDataByYears))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/real_gdp/stats"), app.requirePermissions(economicPermission, app.realGdpDataByYearsStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/real_gdp_per_capita"), app.requirePermissions(economicPermission, app.realGdpPerCapitaDataByYears))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/real_gdp_per_capita/stats"), app.requirePermissions(economicPermission, app.realGdpPerCapitaStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/consumer_sentiment"), app.requirePermissions(economicPermission, app.consumerSentimentDataByYears))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/consumer_sentiment/stats"), app.requirePermissions(economicPermission, app.consumerSentimentStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/retail_sales"), app.requirePermissions(economicPermission, app.retailSalesDataByYears))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/retail_sales/stats"), app.requirePermissions(economicPermission, app.retailSalesDataByYearsStats))

	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/treasury_yield/:maturity"), app.requirePermissions(economicPermission, app.treasuryYieldByYears))
	router.HandlerFunc(http.MethodGet, WithVersion("/%s/economic/treasury_yield/:maturity/stats"), app.requirePermissions(economicPermission, app.treasuryYieldByYearsStats))

	router.HandlerFunc(http.MethodPost, WithVersion("/%s/users"), app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, WithVersion("/%s/users/activated"), app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, WithVersion("/%s/tokens/activation"), app.createActivationTokenHandler)
	router.HandlerFunc(http.MethodPost, WithVersion("/%s/tokens/authentication"), app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}

func WithVersion(pathFmt string) string {
	return fmt.Sprintf(pathFmt, apiVersion)
}
