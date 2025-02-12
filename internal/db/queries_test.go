package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jingen11/stonk-tracker/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertStockPrice(t *testing.T) {
	url := os.Getenv("MONGODB_URL_TEST")
	dbClient, _ := Init(url)
	defer Disconnect(dbClient)

	stonkDb := dbClient.Database("stonk-test")
	priceColl, _ := InitPriceCollection(stonkDb)
	symbolColl, _ := InitSymbolCollection(stonkDb)

	defer priceColl.Drop(context.Background())
	defer symbolColl.Drop(context.Background())

	q := Query{
		SymbolColl: symbolColl,
		PriceColl:  priceColl,
	}

	cases := []struct {
		input      models.StockData
		latestDate string
	}{
		{
			input: models.StockData{
				Status:     "OK",
				From:       "2025-02-03",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			},
			latestDate: "2025-02-03",
		},
		{
			input: models.StockData{
				Status:     "OK",
				From:       "2025-02-10",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			},
			latestDate: "2025-02-10",
		}, {
			input: models.StockData{
				Status:     "OK",
				From:       "2025-02-08",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			},
			latestDate: "2025-02-10",
		},
	}

	for _, c := range cases {
		_, err := q.InsertStockPrice(c.input, context.Background())
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		symbol := models.Symbol{}
		ibm := symbolColl.FindOne(context.Background(), bson.M{"symbol": "AAPL"})
		ibm.Decode(&symbol)
		if symbol.LastFetchedDate.Time().Format("2006-01-02") != c.latestDate {
			t.Fatalf("expected latest date: %s, got date: %s", c.latestDate, symbol.LastFetchedDate.Time().Format("2006-01-02"))
		}
	}
}

func TestInsertSymbolStockPrices(t *testing.T) {
	url := os.Getenv("MONGODB_URL_TEST")
	dbClient, _ := Init(url)
	defer Disconnect(dbClient)

	stonkDb := dbClient.Database("stonk-test")
	priceColl, _ := InitPriceCollection(stonkDb)
	symbolColl, _ := InitSymbolCollection(stonkDb)

	defer priceColl.Drop(context.Background())
	defer symbolColl.Drop(context.Background())

	q := Query{
		SymbolColl: symbolColl,
		PriceColl:  priceColl,
	}

	cases := []struct {
		input        []models.StockData
		symbol       string
		updateCounts int
		latestDate   string
	}{
		{
			input: []models.StockData{{
				Status:     "OK",
				From:       "2025-02-03",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-02-02",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-02-01",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}},
			updateCounts: 3,
			symbol:       "AAPL",
			latestDate:   "2025-02-03",
		},
		{
			input: []models.StockData{{
				Status:     "OK",
				From:       "2025-01-30",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-01-29",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-01-28",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}},
			symbol:       "AAPL",
			updateCounts: 3,
			latestDate:   "2025-02-03",
		},
		{
			input: []models.StockData{{
				Status:     "OK",
				From:       "2025-02-10",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-02-09",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-02-08",
				Symbol:     "AAPL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}},
			updateCounts: 3,
			symbol:       "AAPL",
			latestDate:   "2025-02-10",
		},
		{
			input: []models.StockData{{
				Status:     "OK",
				From:       "2025-02-10",
				Symbol:     "GOOGL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-02-03",
				Symbol:     "GOOGL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}, {
				Status:     "OK",
				From:       "2025-02-02",
				Symbol:     "GOOGL",
				Open:       229.57,
				High:       230.585,
				Low:        227.2,
				Close:      227.65,
				Volume:     30219759,
				AfterHours: 226.9999,
				PreMarket:  228.5,
			}},
			updateCounts: 3,
			symbol:       "GOOGL",
			latestDate:   "2025-02-10",
		},
	}

	for _, c := range cases {
		_, err := q.InsertSymbolStockPrices(c.input, c.symbol, context.Background())
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		symbol := models.Symbol{}
		ibm := symbolColl.FindOne(context.Background(), bson.M{"symbol": "AAPL"})
		ibm.Decode(&symbol)
		if symbol.LastFetchedDate.Time().Format("2006-01-02") != c.latestDate {
			t.Fatalf("expected latest date: %s, got date: %s", c.latestDate, symbol.LastFetchedDate.Time().Format("2006-01-02"))
		}
	}
}

func TestGetAllSymbols(t *testing.T) {
	url := os.Getenv("MONGODB_URL_TEST")
	dbClient, _ := Init(url)
	defer Disconnect(dbClient)

	stonkDb := dbClient.Database("stonk-test")
	symbolColl, _ := InitSymbolCollection(stonkDb)

	defer symbolColl.Drop(context.Background())

	docs := make([]interface{}, 0)
	docs = append(docs, models.Symbol{
		Symbol:          "IBM",
		LastFetchedDate: primitive.NewDateTimeFromTime(time.Now())})
	docs = append(docs, models.Symbol{
		Symbol:          "AAPL",
		LastFetchedDate: primitive.NewDateTimeFromTime(time.Now())})
	docs = append(docs, models.Symbol{
		Symbol:          "GOOGL",
		LastFetchedDate: primitive.NewDateTimeFromTime(time.Now())})

	symbolColl.InsertMany(context.Background(), docs)

	q := Query{
		SymbolColl: symbolColl,
	}

	symbols, err := q.GetAllSymbols(context.Background())

	if err != nil {
		t.Fatalf("Error while getting all symbols: %v", err)
	}

	if len(symbols) != 3 {
		t.Fatalf("Expected length of symbols: %d, Actual length: %d", 3, len(symbols))
	}
}

func TestGetStockPrices(t *testing.T) {
	url := os.Getenv("MONGODB_URL_TEST")
	dbClient, _ := Init(url)
	defer Disconnect(dbClient)

	stonkDb := dbClient.Database("stonk-test")
	priceColl, _ := InitPriceCollection(stonkDb)

	defer priceColl.Drop(context.Background())

	docs := make([]interface{}, 0)
	dateOne := time.Now()
	dateTwo, _ := time.Parse("2006-01-02", "2025-02-10")
	dateThree, _ := time.Parse("2006-01-02", "2025-02-09")
	docs = append(docs, models.Price{
		Symbol: "IBM",
		Open:   255.2800,
		High:   256.9300,
		Low:    252.0200,
		Close:  252.3400,
		Volume: 3370284,
		Date:   primitive.NewDateTimeFromTime(dateOne)})
	docs = append(docs, models.Price{
		Symbol: "IBM",
		Open:   255.2800,
		High:   256.9300,
		Low:    252.0200,
		Close:  252.3400,
		Volume: 3370284,
		Date:   primitive.NewDateTimeFromTime(dateTwo)})
	docs = append(docs, models.Price{
		Symbol: "IBM",
		Open:   255.2800,
		High:   256.9300,
		Low:    252.0200,
		Close:  252.3400,
		Volume: 3370284,
		Date:   primitive.NewDateTimeFromTime(dateThree)})
	docs = append(docs, models.Price{
		Symbol: "GOOGL",
		Open:   255.2800,
		High:   256.9300,
		Low:    252.0200,
		Close:  252.3400,
		Volume: 3370284,
		Date:   primitive.NewDateTimeFromTime(dateOne)})

	priceColl.InsertMany(context.Background(), docs)

	q := Query{
		PriceColl: priceColl,
	}

	cases := []struct {
		input          GetStockPriceOpt
		expectedLength int
	}{
		{input: GetStockPriceOpt{
			Symbol: "IBM",
			limit:  2,
		}, expectedLength: 2},
		{input: GetStockPriceOpt{
			Symbol: "IBM",
			limit:  5,
		}, expectedLength: 3},
		{input: GetStockPriceOpt{
			Symbol: "GOOGL",
			limit:  10,
		}, expectedLength: 1},
	}

	for _, c := range cases {
		prices, err := q.GetStockPrices(context.Background(), &c.input)

		if err != nil {
			t.Fatalf("Error while getting all symbols: %v", err)
		}

		if len(prices) != c.expectedLength {
			t.Fatalf("Expected length of prices: %d, Actual length: %d", c.expectedLength, len(prices))
		}
	}
}
