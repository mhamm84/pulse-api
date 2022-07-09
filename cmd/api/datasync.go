package main

func (app *application) startDataSyncs() {
	app.services.alphaVantageEconomicService.StartDataSyncTask()
}
