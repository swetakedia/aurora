package aurora

func initBlocksafeCoreInfo(app *App) {
	app.UpdateBlocksafeCoreInfo()
}

func init() {
	appInit.Add("blocksafeCoreInfo", initBlocksafeCoreInfo, "app-context", "log")
}
