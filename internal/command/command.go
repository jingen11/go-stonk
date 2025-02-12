package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jingen11/stonk-tracker/internal/models"
	"github.com/jingen11/stonk-tracker/internal/utils"
)

type Command struct {
	Cfg   utils.ProjectConfig
	Input []string
}

func HandleRefresh(p *Command) error {
	symbols, err := p.Cfg.Query.GetAllSymbols(context.TODO())
	if err != nil {
		return err
	}
	type stonkPackage struct {
		symbol string
		dates  []time.Time
	}

	packages := []stonkPackage{}
	for _, symbol := range symbols {
		dates := []time.Time{}
		prevDate := truncateToDay(time.Now().Add(-time.Hour * 24))

		shouldContinue := prevDate.After(symbol.LastFetchedDate.Time())
		for shouldContinue {
			if prevDate.Weekday() != time.Sunday && prevDate.Weekday() != time.Saturday {
				dates = append(dates, prevDate)
			}
			prevDate = prevDate.Add(-time.Hour * 24)
			shouldContinue = prevDate.After(symbol.LastFetchedDate.Time())
		}

		if len(dates) > 0 {
			packages = append(packages, stonkPackage{
				symbol: symbol.Symbol,
				dates:  dates,
			})
		}

		// stockData, err := p.Cfg.ApiClient.GetPrices(symbol.Symbol, time.Now().Format("2006-01-02"))
		// if err != nil {
		// 	fmt.Printf("Error getting stockData for symbol: %s", symbol.Symbol)
		// 	continue
		// }
		// _, err = p.Cfg.Query.InsertStockPrice(stockData, context.TODO())
		// if err != nil {
		// 	fmt.Printf("Error insertting stock price for symbol: %s", symbol.Symbol)
		// 	continue
		// }
	}

	stockChan := make(chan models.StockData)
	errChan := make(chan error)
	total := 0

	for i := 0; i < len(packages); i++ {
		for j := 0; j < len(packages[i].dates); j++ {
			total++
			go getPriceConcurrently(errChan, stockChan, packages[i].symbol, packages[i].dates[j], p)
		}
	}

	stockRes := getPriceChanSubscriber(errChan, stockChan, total)

	type stonkStonksResponse struct {
		symbol    string
		stockData []models.StockData
	}

	mapper := map[string]stonkStonksResponse{}

	for _, stock := range *stockRes {
		if entry, ok := mapper[stock.Symbol]; !ok {
			mapper[stock.Symbol] = stonkStonksResponse{
				symbol: stock.Symbol,
				stockData: []models.StockData{
					stock,
				},
			}
		} else {
			latest := append(entry.stockData, stock)
			entry.stockData = latest
		}
	}

	for k, v := range mapper {
		_, err = p.Cfg.Query.InsertSymbolStockPrices(v.stockData, k, context.TODO())
		if err != nil {
			fmt.Printf("Error inserting stock price for symbol: %s\n", k)
		}
	}

	return nil
}

func HandlerAddNewSymbol(p *Command) error {
	if len(p.Input) != 1 {
		return errors.New("Please provide a stonk symbol")
	}
	symbol := p.Input[0]
	dates := []time.Time{}
	prevDate := time.Now().Add(-time.Hour * 24)

	for len(dates) != p.Cfg.HistoricalTimeFrame {
		if prevDate.Weekday() != time.Sunday && prevDate.Weekday() != time.Saturday {
			dates = append(dates, prevDate)
		}
		prevDate = prevDate.Add(-time.Hour * 24)
	}

	stockChan := make(chan models.StockData)
	errChan := make(chan error)

	for i := 0; i < len(dates); i++ {
		go getPriceConcurrently(errChan, stockChan, symbol, dates[i], p)
	}

	stocks := getPriceChanSubscriber(errChan, stockChan, p.Cfg.HistoricalTimeFrame)

	_, err := p.Cfg.Query.InsertSymbolStockPrices(*stocks, symbol, context.TODO())
	if err != nil {
		fmt.Printf("Error inserting stock price for symbol: %s\n", symbol)
		return err
	}
	return nil
}

// func HandleGetTrendReversal(p *Command) error {
// 	if len(p.input) != 1 {
// 		return errors.New("Please provide a stonk symbol")
// 	}
// 	symbol := p.input[0]
// 	prices, err := p.cfg.Query.GetStockPrices(context.Background(), &db.GetStockPriceOpt{
// 		Symbol: symbol,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	for _, price := range prices {

// 	}
// }

func getPriceConcurrently(errChan chan error, stockChan chan models.StockData, symbol string, date time.Time, p *Command) {
	stockData, err := p.Cfg.ApiClient.GetPrices(symbol, date.Format("2006-01-02"))
	if err != nil { // holiday will cause error, so it is ok to swallow
		errChan <- err
		return
	}
	stockChan <- stockData
}

func getPriceChanSubscriber(errChan chan error, stockChan chan models.StockData, length int) *[]models.StockData {
	stocks := []models.StockData{}
	count := 0
	ended := false
	for !ended {
		select {
		case err := <-errChan:
			count++
			fmt.Printf("Error fetching stock price for symbol for %d days: %v\n", length, err)
			if count == length {
				ended = true
			}
		case s := <-stockChan:
			count++
			stocks = append(stocks, s)
			if count == length {
				ended = true
			}
		}
	}
	return &stocks
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
