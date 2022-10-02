package api

import (
	"errors"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/services"
	"github.com/mhamm84/pulse-api/internal/validator"
	"net/http"
)

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token string `json:"token"`
	}

	err := app.ReadJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	v := validator.New()
	if services.ValidateTokenPlaintext(v, input.Token); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.services.UserService.ActivateUser(r.Context(), input.Token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.WriteJson(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	user := &data.User{
		Name: input.Name, Email: input.Email, Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.services.UserService.RegisterUser(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write a JSON response containing the user data along with a 201 Created status // code.
	err = app.WriteJson(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
