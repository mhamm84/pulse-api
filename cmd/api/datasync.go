package main

func (app *application) startDataSyncs() {
	app.services.cpiService.StartCpiDataSyncTask()
}
