package cmd

import (
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	customerCfg "tesodev-korpes/CustomerService/config"
	"tesodev-korpes/CustomerService/internal"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
)

// @title Customer Service API
// @version 1.0
// @description API documentation for Customer Service.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8001
// @BasePath /api/v1
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

func BootCustomerService(client *mongo.Client, e *echo.Echo) {
	config := customerCfg.GetCustomerConfig("dev")
	customerCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewRepository(customerCol)
	service := internal.NewService(repo)
	handler := internal.NewHandler(service)
	handler.Route(e)

	e.Logger.Fatal(e.Start(config.Port))
}
