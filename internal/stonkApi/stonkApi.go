package stonkapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jingen11/stonk-tracker/internal/models"
)

const POLYGON_IO_HOST_NAME = "api.polygon.io"

type StonkApiClient struct {
	Client          http.Client
	roundRobinIndex int
	apiKeys         []string
}

type ErrorResponse struct {
	Message string `json:"message"`
	Url     string
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

func (client *StonkApiClient) GetPrices(symbol, date string) (models.StockData, error) {
	endpoint := url.URL{
		Scheme:   "https",
		Host:     POLYGON_IO_HOST_NAME,
		Path:     "v1/open-close/" + symbol + "/" + date,
		RawQuery: fmt.Sprintf("adjusted=true&apiKey=%s", client.roundRobinGetApiKey()),
	}
	res, err := client.Client.Get(endpoint.String())

	stockData := models.StockData{}

	if err != nil {
		return stockData, err
	}

	if res.StatusCode == 429 {
		fmt.Println("Cooling down stonk api")
		time.Sleep(1 * time.Minute)
		res, err = client.Client.Get(endpoint.String())
		if err != nil {
			return stockData, err
		}
	}

	if res.StatusCode > 299 {
		e := ErrorResponse{}
		json.NewDecoder(res.Body).Decode(&e)
		e.Url = endpoint.String()
		return stockData, errors.New(fmt.Sprintf("url: %s,\n message: %s", e.Url, e.Message))
	}

	err = json.NewDecoder(res.Body).Decode(&stockData)
	defer res.Body.Close()

	if err != nil {
		return stockData, err
	}

	return stockData, nil
}
