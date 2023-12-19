package main

import (
	"github.com/wuttinanhi/code-judge-system/cmds"
	"github.com/wuttinanhi/code-judge-system/services"
)

func main() {
	serviceKit := services.CreateServiceKit()
	app := cmds.SetupWeb(serviceKit)
	app.Listen(":3000")
}
