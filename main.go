package main

import (
	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/services"
)

func main() {
	db := databases.NewMySQLDatabase()
	serviceKit := services.CreateServiceKit(db)

	app := controllers.SetupWeb(serviceKit)
	app.Listen(":3000")
}
