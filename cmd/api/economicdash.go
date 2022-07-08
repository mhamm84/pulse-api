package main

import "net/http"

func (app *application) economicDashHandler(w http.ResponseWriter, r *http.Request) {
	data, err := app.services.economicdashservice.GetDashboardSummary()
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorHandler(w, r)
	}
	env := envelope{
		"summaries": &data,
	}
	app.WriteJson(w, http.StatusOK, env, nil)
}
