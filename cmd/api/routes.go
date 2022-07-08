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

	return app.recoverPanic(router)
}

func pathWithVersion(pathFmt string) string {
	return fmt.Sprintf(pathFmt, apiVersion)
}
