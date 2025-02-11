package stonkapi

import (
	"os"
	"testing"
)

func TestRoundRobinGetApiKey(t *testing.T) {
	cases := []struct {
		Inputs          []string
		RoundRobinRound int
		ExpectedKey     string
	}{
		{
			Inputs: []string{
				"1", "2", "3",
			},
			RoundRobinRound: 1,
			ExpectedKey:     "1",
		},
		{
			Inputs: []string{
				"1", "2", "3",
			},
			RoundRobinRound: 2,
			ExpectedKey:     "2",
		},
		{
			Inputs: []string{
				"1", "2", "3",
			},
			RoundRobinRound: 3,
			ExpectedKey:     "3",
		},
		{
			Inputs: []string{
				"1", "2", "3",
			},
			RoundRobinRound: 4,
			ExpectedKey:     "1",
		},
		{
			Inputs: []string{
				"1",
			},
			RoundRobinRound: 1,
			ExpectedKey:     "1",
		},
		{
			Inputs: []string{
				"1",
			},
			RoundRobinRound: 2,
			ExpectedKey:     "1",
		},
	}

	for i, c := range cases {
		client := InitStonkApiClient(c.Inputs)
		apiKey := ""
		for round := 0; round < c.RoundRobinRound; round++ {
			apiKey = client.roundRobinGetApiKey()
		}

		if apiKey != c.ExpectedKey {
			t.Fatalf("Test case %d: expected key: %s, actual key: %s", i, c.ExpectedKey, apiKey)
		}
	}
}

func TestGetPrices(t *testing.T) {
	apiKey := os.Getenv("POLYGON_IO_KEY_1")
	client := InitStonkApiClient([]string{apiKey})
	stonkData, err := client.GetPrices("AAPL", "2025-02-10")
	if err != nil {
		t.Fatalf("error getting price from api, error: %v", err)
	}
	if stonkData.Symbol != "AAPL" {
		t.Fatalf("mismatch stonk symbol")
	}
}

func TestGetPricesRateLimit(t *testing.T) {
	apiKey := os.Getenv("POLYGON_IO_KEY_2")
	client := InitStonkApiClient([]string{apiKey})
	for i := 0; i < 6; i++ {
		_, err := client.GetPrices("AAPL", "2025-02-10")

		if err != nil {
			t.Fatalf("error getting price from api rate limit, error: %v", err)
		}
	}

}
