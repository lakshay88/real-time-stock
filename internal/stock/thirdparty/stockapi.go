package thirdparty

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lakshay88/real-time-stock/config"
)

type StockAPIClient struct {
	BaseUrl    string
	ApiKey     string
	HTTPClient *http.Client
}

func SetUpStockAPI(cfg config.AppConfig) *StockAPIClient {

	clientRequest := &StockAPIClient{
		BaseUrl: cfg.APIConfiguration.URL,
		ApiKey:  cfg.APIConfiguration.URL,
		HTTPClient: &http.Client{
			Timeout: time.Duration(cfg.APIConfiguration.TimeOut) * time.Second,
		},
	}

	return clientRequest
}

func (c *StockAPIClient) GetStockData(symbol string) ([]byte, error) {
	url := fmt.Sprintf("%s?function=TIME_SERIES_INTRADAY&symbol=%s&interval=1min&apikey=%s", c.BaseUrl, symbol, c.ApiKey)

	var body []byte
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return body, fmt.Errorf("error fetching stock data: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err = io.ReadAll(resp.Body) // Use io.ReadAll() to read entire response
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
