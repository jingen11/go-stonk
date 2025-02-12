package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jingen11/stonk-tracker/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	SymbolColl *mongo.Collection
	PriceColl  *mongo.Collection
}

type GetStockPriceOpt struct {
	Symbol string
	limit  int64
}

func (q *Query) InsertStockPrice(stock models.StockData, ctx context.Context) (*models.Price, error) {
	p := models.Price{}
	symbol := stock.Symbol

	symbolDoc := q.SymbolColl.FindOne(ctx, bson.M{"symbol": symbol})

	if symbolDoc.Err() != nil && symbolDoc.Err() != mongo.ErrNoDocuments {
		fmt.Println("Failed to find symbol")
		return &p, symbolDoc.Err()
	}

	symbolStruct := models.Symbol{}

	if symbolDoc.Err() == mongo.ErrNoDocuments {
		lastFetchedDate, err := time.Parse("2006-01-02", stock.From)
		if err != nil {
			fmt.Println("Failed to parse stonkDate")
			return &p, err
		}

		newSymbol := models.Symbol{
			Symbol:          symbol,
			LastFetchedDate: primitive.NewDateTimeFromTime(lastFetchedDate),
		}
		insertedSymbol, err := q.SymbolColl.InsertOne(ctx, newSymbol)

		if err != nil {
			fmt.Println("Failed to insert symbol")
			return &p, err
		}

		newSymbol.Id = insertedSymbol.InsertedID.(primitive.ObjectID)
		symbolStruct = newSymbol
	} else {
		err := symbolDoc.Decode(&symbolStruct)

		if err != nil {
			fmt.Println("Failed to decode symbol")
			return &p, err
		}
	}

	lastFetchedDate := symbolStruct.LastFetchedDate.Time()

	stonkDate, err := time.Parse("2006-01-02", stock.From)

	if err != nil {
		fmt.Println("Failed to parse stonkDate")
		return &p, err
	}

	p = models.Price{
		Symbol: symbol,
		Date:   primitive.NewDateTimeFromTime(stonkDate),
		Open:   stock.Open,
		Close:  stock.Close,
		High:   stock.High,
		Low:    stock.Low,
		Volume: stock.Volume,
	}

	inserted, err := q.PriceColl.InsertOne(ctx, p)

	if err != nil {
		fmt.Println("Failed to parse stonkDate")
		return &p, err
	}

	p.Id = inserted.InsertedID.(primitive.ObjectID)

	if stonkDate.After(lastFetchedDate) {
		_, err := q.SymbolColl.UpdateByID(ctx, symbolStruct.Id, bson.D{
			{"$set", bson.D{{"lastFetchedDate", primitive.NewDateTimeFromTime(stonkDate)}}},
		})

		if err != nil {
			fmt.Println("failed to update symbol lastFetch")
			return &p, err
		}
	}

	return &p, nil
}

func (q *Query) GetAllSymbols(ctx context.Context) ([]models.Symbol, error) {
	symbolCursor, err := q.SymbolColl.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println("failed to get all symbols")
		return nil, err
	}

	var symbols []models.Symbol

	err = symbolCursor.All(ctx, &symbols)

	if err != nil {
		fmt.Println("faile to decode all symbols")
		return nil, err
	}

	return symbols, nil
}

func (q *Query) GetStockPrices(ctx context.Context, opt *GetStockPriceOpt) ([]models.Price, error) {
	if opt.limit == 0 {
		opt.limit = 100
	}
	opts := options.Find()
	opts.SetLimit(opt.limit)
	opts.SetSort(bson.M{
		"date": -1,
	})
	filters := bson.M{
		"date": bson.M{
			"$lte": primitive.NewDateTimeFromTime(time.Now()),
		},
		"symbol": opt.Symbol,
	}
	priceCursor, err := q.PriceColl.Find(ctx, filters, opts)
	if err != nil {
		fmt.Println("failed to get price for symbol")
		return nil, err
	}
	var prices []models.Price

	err = priceCursor.All(ctx, &prices)

	if err != nil {
		fmt.Println("faile to decode all prices")
		return nil, err
	}

	return prices, nil
}

func (q *Query) InsertSymbolStockPrices(stocks []models.StockData, symbol string, ctx context.Context) (int, error) {
	if len(stocks) == 0 {
		return 0, errors.New("no stocks found")
	}
	symbolDoc := q.SymbolColl.FindOne(ctx, bson.M{"symbol": symbol})

	if symbolDoc.Err() != nil && symbolDoc.Err() != mongo.ErrNoDocuments {
		fmt.Println("Failed to find symbol")
		return 0, symbolDoc.Err()
	}

	symbolStruct := models.Symbol{}

	if symbolDoc.Err() == mongo.ErrNoDocuments {
		lastFetchedDate, err := time.Parse("2006-01-02", "1970-01-01")
		if err != nil {
			fmt.Println("Failed to parse lastFetchedDate")
			return 0, err
		}
		newSymbol := models.Symbol{
			Symbol:          symbol,
			LastFetchedDate: primitive.NewDateTimeFromTime(lastFetchedDate),
		}
		insertedSymbol, err := q.SymbolColl.InsertOne(ctx, newSymbol)

		if err != nil {
			fmt.Println("Failed to insert symbol")
			return 0, err
		}

		newSymbol.Id = insertedSymbol.InsertedID.(primitive.ObjectID)
		symbolStruct = newSymbol
	} else {
		err := symbolDoc.Decode(&symbolStruct)

		if err != nil {
			fmt.Println("Failed to decode symbol")
			return 0, err
		}
	}

	docs := make([]interface{}, 0)

	lastFetchedDate := symbolStruct.LastFetchedDate.Time() // 1970-01-01 || lags behind date remains constant

	latestDate := lastFetchedDate // to be updated for each iteration

	for _, stonk := range stocks {
		stonkDate, err := time.Parse("2006-01-02", stonk.From)
		if err != nil {
			fmt.Println("Failed to parse stonkDate")
			return 0, err
		}
		if stonkDate.After(latestDate) {
			latestDate = stonkDate
		}

		docs = append(docs, models.Price{
			Symbol: symbol,
			Date:   primitive.NewDateTimeFromTime(stonkDate),
			Open:   stonk.Open,
			Close:  stonk.Close,
			High:   stonk.High,
			Low:    stonk.Low,
			Volume: stonk.Volume,
		})
	}

	res, err := q.PriceColl.InsertMany(ctx, docs)

	if err != nil {
		fmt.Println("failed to insert prices")
		return 0, err
	}

	if latestDate.After(lastFetchedDate) {
		_, err := q.SymbolColl.UpdateByID(ctx, symbolStruct.Id, bson.D{
			{"$set", bson.D{{"lastFetchedDate", primitive.NewDateTimeFromTime(latestDate)}}},
		})

		if err != nil {
			fmt.Println("failed to update symbol lastFetch")
			return 0, err
		}
	}

	return len(res.InsertedIDs), nil
}
