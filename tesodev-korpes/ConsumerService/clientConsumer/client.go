package clientConsumer

import (
	"fmt"
	"tesodev-korpes/ConsumerService/internal/types"
	"tesodev-korpes/pkg"
)

type ConsumerClient struct {
	RestClient *pkg.RestClient
}

func NewConsumerClient(restClient *pkg.RestClient) *ConsumerClient {
	return &ConsumerClient{
		RestClient: restClient,
	}
}

func (c *ConsumerClient) GetOrder(orderID string, token string) (*types.OrderResponseModel, error) {
	var res types.OrderResponseModel
	uri := fmt.Sprintf("http://order-service:1881/order/%s", orderID)
	err := c.RestClient.DoGetRequest(uri, &res, token)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *ConsumerClient) GetCustomer(customerID string, token string) (*types.CustomerResponseModel, error) {
	var res types.CustomerResponseModel
	uri := fmt.Sprintf("http://customer-service:1907/customer/%s", customerID)
	err := c.RestClient.DoGetRequest(uri, &res, token)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
