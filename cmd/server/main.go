package main

import (
	router "motorbike-rental-backend/api/routes"
	"motorbike-rental-backend/pkg/app"
	"os"
)

var Version = "v1.0.0"
var BuildTime = "Bilinmiyor"

func main() {
	r := router.NewIdareRouter()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		performMigration(*r)
		return
	}

	a := app.New(r, Version, BuildTime)
	a.Start()
}

func performMigration(r router.IdareRouter) {
	a := app.New(r, Version, BuildTime)
	a.MigrateDB()
}
