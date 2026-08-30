package httpclient

import (
	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
)

type CustomerClient struct {
	*BaseClient
}

func NewCustomerClient(baseURL, token string) *CustomerClient {
	return &CustomerClient{BaseClient: NewBaseClient(baseURL, token)}
}

func (c *CustomerClient) Get(customerID string) (*dto.CustomerResponse, error) {
	var resp dto.CustomerResponse
	if err := c.do("GET", "/api/v1/customers/"+customerID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
