package client

import (
	"fmt"
	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg"
)

type CustomerClient struct {
	RestClient *pkg.RestClient
}

func NewCustomerClient(restClient *pkg.RestClient) *CustomerClient {
	return &CustomerClient{
		RestClient: restClient,
	}
}

func (c *CustomerClient) GetCustomerByID(customerID string, token string) (*types.CustomerResponse, error) {
	var res types.CustomerResponse
	uri := fmt.Sprintf("http://customer-service:1907/customer/%s", customerID)
	err := c.RestClient.DoGetRequest(uri, &res, token)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
