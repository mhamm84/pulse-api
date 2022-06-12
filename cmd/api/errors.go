package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Printf(err.Error())
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.WriteJson(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(status)
	}
}

func (app *application) serverErrorHandler(w http.ResponseWriter, r *http.Request) {
	msg := "The server could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	msg := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestHandler(w http.ResponseWriter, r *http.Request) {
	message := "Bad Request was sent to the server"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}
