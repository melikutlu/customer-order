package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"tesodev-korpes/ConsumerService/clientConsumer"
	"tesodev-korpes/ConsumerService/config"
	"tesodev-korpes/ConsumerService/internal/types"
	"tesodev-korpes/pkg/Kafka/consumer"
	"time"
)

var secretKey string

func init() {

	appConf := config.GetAppConfig("dev")
	secretKey = appConf.SecretKey

}

type Service struct {
	repo          *FinanceRepository
	kafkaConsumer *consumer.Consumer
	conClient     *clientConsumer.ConsumerClient
}

func NewService(repo *FinanceRepository, conClient *clientConsumer.ConsumerClient, kafkaConsumer *consumer.Consumer, brokers []string, topic string) *Service {

	kafkaConsumer.Topic = topic
	kafkaConsumer.CreateConnection(brokers)

	return &Service{
		repo:          repo,
		conClient:     conClient,
		kafkaConsumer: kafkaConsumer,
	}
}

func (s *Service) Read() error {

	ctx := context.Background()
	consumerAction := func(msg string, err error) {
		if err != nil {
			fmt.Printf("Error consuming message: %v\n", err)
			return
		}
		fmt.Printf("Consumed message: %s\n", msg)

		err = s.aggregateCustomerOrder(ctx, msg)
		if err != nil {
			fmt.Printf("Error processing message: %v\n", err)
		}
	}

	go s.kafkaConsumer.Read(consumerAction)

	return nil
}

func (s *Service) aggregateCustomerOrder(ctx context.Context, msg string) error {

	order, err := s.conClient.GetOrder(msg, secretKey)
	if err != nil {
		log.Errorf("Error getting order by ID: %v", err)
		return nil
	}
	if order == nil {
		log.Errorf("Order not found for ID: %s", msg)
		return nil
	}
	log.Infof("Order Info: %+v", order)

	customer, err := s.conClient.GetCustomer(order.CustomerId, secretKey)
	if err != nil {
		log.Errorf("Error getting customer by ID: %v", err)
		return nil
	}
	if customer == nil {
		log.Errorf("Customer not found for ID: %s", order.CustomerId)
		return nil
	}
	log.Infof("Customer Info: %+v", customer)

	priceWithVat := CalculateVat(order.Price)
	log.Infof("Price with VAT: %d", priceWithVat)

	order.Price = priceWithVat

	aggregateData := &types.AggregateData{
		Id:             uuid.New().String(),
		FirstName:      customer.FirstName,
		LastName:       customer.LastName,
		Username:       customer.Username,
		CustomerId:     order.CustomerId,
		OrderName:      order.OrderName,
		ShipmentStatus: order.ShipmentStatus,
		PaymentMethod:  order.PaymentMethod,
		OrderTotal:     order.OrderTotal,
		Price:          order.Price,
		OrderCreatedAt: time.Now(),
		OrderUpdatedAt: time.Now(),
		CreatedAt:      time.Now(),
	}

	_, err = s.repo.Create(ctx, aggregateData)
	if err != nil {
		log.Errorf("Error saving consumer to repository: %v", err)
		return nil
	}

	log.Infof("CustomerOrder saved successfully: %+v", aggregateData)
	return nil
}
