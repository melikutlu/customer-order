package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"os/signal"
	"syscall"
	"tesodev-korpes/ConsumerService/clientConsumer"
	consumerCmd "tesodev-korpes/ConsumerService/cmd"
	"tesodev-korpes/CustomerService/cmd"
	"tesodev-korpes/OrderService/client"
	_ "tesodev-korpes/OrderService/client"
	orderCmd "tesodev-korpes/OrderService/cmd"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/Kafka/consumer"
	"tesodev-korpes/pkg/Kafka/producer"
	"tesodev-korpes/pkg/log"
	"tesodev-korpes/shared/config"
)

func main() {

	dbConf := config.GetDBConfig("dev")

	mongoClient, err := pkg.GetMongoClient(dbConf.MongoDuration, dbConf.MongoClientURI)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	hClient := client.NewCustomerClient(pkg.NewRestClient())
	consumerClient := clientConsumer.NewConsumerClient(pkg.NewRestClient())

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(pkg.CorrelationIDMiddleware)
	e.Use(log.Logger())
	e.Use(pkg.Authenticate)

	brokers := []string{"kafka:9092"}
	topic := "product_price_change_events"

	kafkaProducer := producer.NewProducer(brokers, topic)

	kafkaConsumer := &consumer.Consumer{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if len(os.Args) < 2 {
		panic("Please provide a service to start: customer, order, or both")
	}
	input := os.Args[1]

	switch input {
	case "customer":
		cmd.BootCustomerService(mongoClient, e)
	case "order":
		orderCmd.BootOrderService(mongoClient, hClient, kafkaProducer, e)
	case "consumer":
		go consumerCmd.BootConsumerService(mongoClient, kafkaConsumer, consumerClient, e, brokers, topic)
	case "both":

		go cmd.BootCustomerService(mongoClient, e)
		go orderCmd.BootOrderService(mongoClient, hClient, kafkaProducer, e)
		go consumerCmd.BootConsumerService(mongoClient, kafkaConsumer, consumerClient, e, brokers, topic)
	default:
		panic("Invalid input. Use 'customer', 'order', or 'both'.")
	}

	<-sigs
	fmt.Println("Shutting down...")

	kafkaConsumer.Close()
	kafkaProducer.Close()

	fmt.Println("Kafka connections closed. Exiting.")

	i := 0
	for i < 5 {
		i++
	}
}
