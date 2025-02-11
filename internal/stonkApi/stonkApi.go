package stonkapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jingen11/stonk-tracker/internal/models"
)

const ALPHA_VANTAGE_HOST_NAME = "www.alphavantage.co/query"

type StonkApiClient struct {
	Client          http.Client
	roundRobinIndex int
	apiKeys         []string
}

func InitStonkApiClient(apiKeys []string) *StonkApiClient {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	return &StonkApiClient{
		Client:          client,
		roundRobinIndex: 0,
		apiKeys:         apiKeys,
	}
}

func (client *StonkApiClient) roundRobinGetApiKey() string {
	curIndex := client.roundRobinIndex

	apiKey := client.apiKeys[curIndex]

	if curIndex == len(client.apiKeys)-1 {
		client.roundRobinIndex = 0
	} else {
		client.roundRobinIndex = curIndex + 1
	}

	return apiKey
}

func (client *StonkApiClient) GetPrices(symbol string) (models.StockData, error) {
	endpoint := url.URL{
		Scheme:   "https",
		Host:     ALPHA_VANTAGE_HOST_NAME,
		RawQuery: fmt.Sprintf("function=%s&symbol=%s&apikey=%s", "TIME_SERIES_DAILY_ADJUSTED", symbol, client.roundRobinGetApiKey()),
	}
	res, err := client.Client.Get(endpoint.String())

	stockData := models.StockData{}

	if err != nil {
		return stockData, err
	}

	err = json.NewDecoder(res.Body).Decode(&stockData)
	defer res.Body.Close()

	if err != nil {
		return stockData, err
	}

	return stockData, nil
}
