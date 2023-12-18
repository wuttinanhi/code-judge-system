package main

import (
	"github.com/wuttinanhi/code-judge-system/cmds"
	"github.com/wuttinanhi/code-judge-system/services"
)

func main() {
	services.InitServiceKit()
	app := cmds.SetupWeb()
	app.Listen(":3000")
}
