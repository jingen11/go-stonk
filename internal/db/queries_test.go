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

func TestInsertStockPrices(t *testing.T) {
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
		input        models.StockData
		updatedCount int
		latestDate   string
	}{
		{
			input: models.StockData{
				MetaData: models.Meta{
					Information:   "Daily Prices (open, high, low, close) and Volumes",
					Symbol:        "IBM",
					LastRefreshed: "2025-02-07",
					OutputSize:    "Full size",
					TimeZone:      "US/Eastern",
				},
				TimeSeriesDaily: map[string]models.DailyPrice{
					"2025-02-03": {
						Open:   "255.2800",
						High:   "256.9300",
						Low:    "252.0200",
						Close:  "252.3400",
						Volume: "3370284",
					},
					"2025-02-02": {
						Open:   "262.9800",
						High:   "263.3800",
						Low:    "252.7300",
						Close:  "253.4400",
						Volume: "6128293",
					},
					"2025-02-01": {
						Open:   "265.7100",
						High:   "265.7200",
						Low:    "261.1800",
						Close:  "263.3000",
						Volume: "6165096",
					},
				},
			},
			updatedCount: 3,
			latestDate:   "2025-02-03",
		},
		{
			input: models.StockData{
				MetaData: models.Meta{
					Information:   "Daily Prices (open, high, low, close) and Volumes",
					Symbol:        "IBM",
					LastRefreshed: "2025-02-07",
					OutputSize:    "Full size",
					TimeZone:      "US/Eastern",
				},
				TimeSeriesDaily: map[string]models.DailyPrice{
					"2025-02-07": {
						Open:   "255.2800",
						High:   "256.9300",
						Low:    "252.0200",
						Close:  "252.3400",
						Volume: "3370284",
					},
					"2025-02-06": {
						Open:   "262.9800",
						High:   "263.3800",
						Low:    "252.7300",
						Close:  "253.4400",
						Volume: "6128293",
					},
					"2025-02-05": {
						Open:   "265.7100",
						High:   "265.7200",
						Low:    "261.1800",
						Close:  "263.3000",
						Volume: "6165096",
					},
					"2025-02-04": {
						Open:   "260.0000",
						High:   "265.2500",
						Low:    "258.1233",
						Close:  "264.4600",
						Volume: "6077652",
					},
				},
			},
			updatedCount: 4,
			latestDate:   "2025-02-07",
		},
		{
			input: models.StockData{
				MetaData: models.Meta{
					Information:   "Daily Prices (open, high, low, close) and Volumes",
					Symbol:        "IBM",
					LastRefreshed: "2025-02-07",
					OutputSize:    "Full size",
					TimeZone:      "US/Eastern",
				},
				TimeSeriesDaily: map[string]models.DailyPrice{
					"2025-02-07": {
						Open:   "255.2800",
						High:   "256.9300",
						Low:    "252.0200",
						Close:  "252.3400",
						Volume: "3370284",
					},
					"2025-02-06": {
						Open:   "262.9800",
						High:   "263.3800",
						Low:    "252.7300",
						Close:  "253.4400",
						Volume: "6128293",
					},
					"2025-02-05": {
						Open:   "265.7100",
						High:   "265.7200",
						Low:    "261.1800",
						Close:  "263.3000",
						Volume: "6165096",
					},
					"2025-02-04": {
						Open:   "260.0000",
						High:   "265.2500",
						Low:    "258.1233",
						Close:  "264.4600",
						Volume: "6077652",
					},
				},
			},
			updatedCount: 0,
			latestDate:   "2025-02-07",
		},
	}

	for _, c := range cases {
		count, err := q.InsertStockPrices(c.input, context.Background())
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if count != c.updatedCount {
			t.Fatalf("expected updated count: %d, got count: %d", c.updatedCount, count)
		}
		symbol := models.Symbol{}
		ibm := symbolColl.FindOne(context.Background(), bson.M{"symbol": "IBM"})
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
		Open:   "255.2800",
		High:   "256.9300",
		Low:    "252.0200",
		Close:  "252.3400",
		Volume: "3370284",
		Date:   primitive.NewDateTimeFromTime(dateOne)})
	docs = append(docs, models.Price{
		Symbol: "IBM",
		Open:   "255.2800",
		High:   "256.9300",
		Low:    "252.0200",
		Close:  "252.3400",
		Volume: "3370284",
		Date:   primitive.NewDateTimeFromTime(dateTwo)})
	docs = append(docs, models.Price{
		Symbol: "IBM",
		Open:   "255.2800",
		High:   "256.9300",
		Low:    "252.0200",
		Close:  "252.3400",
		Volume: "3370284",
		Date:   primitive.NewDateTimeFromTime(dateThree)})
	docs = append(docs, models.Price{
		Symbol: "GOOGL",
		Open:   "255.2800",
		High:   "256.9300",
		Low:    "252.0200",
		Close:  "252.3400",
		Volume: "3370284",
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
