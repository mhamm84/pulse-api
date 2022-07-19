package main

func (app *application) startEconomicReportDataSync() {

	app.services.alphaVantageEconomicService.StartDataSyncTask()
}
