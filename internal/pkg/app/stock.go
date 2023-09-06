package app

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient interface used to mock http request to test
type HTTPClient interface {
	// Do ...
	Do(req *http.Request) (*http.Response, error)
}

// StockApp ...
type StockApp struct {
	stockURL string
	client   HTTPClient
}

// NewStockApp creates new implementation of Stock service
func NewStockApp(
	stockURL string,
	client HTTPClient,
) *StockApp {
	return &StockApp{
		stockURL: stockURL,
		client:   client,
	}
}

// HandleStockRequest ...
func (s *StockApp) HandleStockRequest(key string) string {
	return s.getStockPrice(key)
}

func (s *StockApp) getStockPrice(key string) string {
	stockServiceURL := fmt.Sprintf(s.stockURL, url.QueryEscape(key))
	log.Println("processing", stockServiceURL)

	request, err := http.NewRequest(http.MethodGet, stockServiceURL, nil)

	response, err := s.client.Do(request)
	if err != nil {
		log.Println("error :", err)
		return "stock service unreachable"
	}

	if response.StatusCode == http.StatusOK {
		content, err := csv.NewReader(response.Body).ReadAll()
		if err != nil {
			log.Println("error :", err)
			return "could not parse CSV"
		}
		symbol := content[1][0]
		price := content[1][6]
		log.Println("content:", content)
		if price == "N/D" {
			return fmt.Sprintf("%s quote is not available", key)
		}
		return fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(symbol), price)
	}

	log.Println("error : response.StatusCode is ", response.StatusCode)
	return "stock service is not available"
}
