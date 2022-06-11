package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "up",
		"environment": app.cfg.env,
		"version":     version,
	}

	err := app.WriteJson(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Printf(err.Error())
		http.Error(w, "The server could not process the request", http.StatusInternalServerError)
	}
}
