package main

import (
	"os"

	"github.com/wuttinanhi/code-judge-system/consumers"
	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/services"
)

func main() {
	db := databases.NewMySQLDatabase()
	serviceKit := services.CreateServiceKit(db)

	MODE := os.Getenv("MODE")

	if MODE == "CONSUMER" {
		consumers.StartSubmissionConsumer(serviceKit)
		return
	}

	app := controllers.SetupWeb(serviceKit)
	app.Listen(":3000")
}
