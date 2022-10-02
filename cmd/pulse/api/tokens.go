package api

import (
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/validator"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := app.ReadJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.services.UserService.GetByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email found.")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if user.Activated {
		v.AddError("user", "user already activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	token, err := app.services.TokenService.New(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID}
		err = app.mailer.Send(user.Email, "token_activation.tmpl", data)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	// Validate the email and password provided by the client.
	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.services.UserService.GetByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Match(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.services.TokenService.New(r.Context(), user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.WriteJson(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
