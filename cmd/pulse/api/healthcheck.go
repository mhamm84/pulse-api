package api

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "up",
		"system_info": map[string]string{
			"environment": app.cfg.Env,
			"version":     "1.0.0",
		},
	}

	err := app.WriteJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
