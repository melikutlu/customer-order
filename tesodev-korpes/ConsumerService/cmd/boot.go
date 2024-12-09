package cmd

import (
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/ConsumerService/clientConsumer"
	config4 "tesodev-korpes/ConsumerService/config"
	"tesodev-korpes/ConsumerService/internal"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/Kafka/consumer"
)

func BootConsumerService(client *mongo.Client, kafkaConsumer *consumer.Consumer, conClient *clientConsumer.ConsumerClient, e *echo.Echo, brokers []string, topic string) {

	config := config4.GetConsumerConfig("dev")

	consumerCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewFinanceRepository(consumerCol)
	service := internal.NewService(repo, conClient, kafkaConsumer, brokers, topic)
	service.Read()

	e.Logger.Fatal(e.Start(config.Port))
}
