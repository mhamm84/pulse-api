package main

import "net/http"

func (app *application) economicDashHandler(w http.ResponseWriter, r *http.Request) {
	data, err := app.services.economicdashservice.GetDashboardSummary()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	env := envelope{
		"economicSummaries": &data,
	}
	app.WriteJson(w, http.StatusOK, env, nil)
}
