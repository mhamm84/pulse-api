package main

import (
	"net/http"
	"time"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "up",
		"system_info": map[string]string{
			"environment": app.cfg.env,
			"version":     version,
		},
	}

	time.Sleep(5 * time.Second)

	err := app.WriteJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Printf(err.Error())
		app.serverErrorHandler(w, r)
	}
}
