package main

import (
	"github.com/spf13/viper"
	"github.com/wuttinanhi/code-judge-system/configs"
	"github.com/wuttinanhi/code-judge-system/consumers"
	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/services"
)

func main() {
	configs.LoadConfig()

	db := databases.NewMySQLDatabase()
	serviceKit := services.CreateServiceKit(db)

	APP_MODE := viper.GetString("APP_MODE")

	if APP_MODE == "CONSUMER" {
		consumers.StartSubmissionConsumer(serviceKit)
		return
	}

	api := controllers.SetupAPI(serviceKit)
	api.Listen(":3000")
}
