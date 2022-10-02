package api

import (
	"fmt"
	"github.com/mhamm84/pulse-api/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	utils.Logger(r.Context()).Error("application error",
		zap.Error(err),
		zap.Any("request", &r),
	)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.WriteJson(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(status)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	msg := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request) {
	message := "Bad Request was sent to the server"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "Too many requests were made"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
