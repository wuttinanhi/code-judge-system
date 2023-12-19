package main

import (
	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/services"
)

func main() {
	serviceKit := services.CreateServiceKit()
	app := controllers.SetupWeb(serviceKit)
	app.Listen(":3000")
}
