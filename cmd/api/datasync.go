package main

func (app *application) startDataSyncs() {
	app.services.cpiService.StartDataSyncTask()
	app.services.consumerSentimentService.StartDataSyncTask()
}
