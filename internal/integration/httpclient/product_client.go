package httpclient

import (
	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
)

type ProductClient struct {
	*BaseClient
}

func NewProductClient(baseURL, token string) *ProductClient {
	return &ProductClient{BaseClient: NewBaseClient(baseURL, token)}
}

func (c *ProductClient) Get(productID string) (*dto.ProductResponse, error) {
	var resp dto.ProductResponse
	if err := c.do("GET", "/api/v1/products/"+productID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *ProductClient) DecreaseStock(productID string, quantity int) error {
	return c.do("PATCH", "/api/v1/products/"+productID+"/decrease-stock",
		dto.AdjustStockRequest{Quantity: quantity}, nil)
}

func (c *ProductClient) IncreaseStock(productID string, quantity int) error {
	return c.do("PATCH", "/api/v1/products/"+productID+"/increase-stock",
		dto.AdjustStockRequest{Quantity: quantity}, nil)
}
