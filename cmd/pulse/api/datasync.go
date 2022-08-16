package api

func (app *application) startEconomicReportDataSync() {

	app.services.AlphaVantageEconomicService.StartDataSyncTask()
}
