package api

func (app *application) startEconomicReportDataSync() {

	app.services.alphaVantageEconomicService.StartDataSyncTask()
}
