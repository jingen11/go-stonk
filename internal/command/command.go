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
	for _, symbol := range symbols {
		stockData, err := p.Cfg.ApiClient.GetPrices(symbol.Symbol, time.Now().Format("2006-01-02"))
		if err != nil {
			fmt.Printf("Error getting stockData for symbol: %s", symbol.Symbol)
			continue
		}
		_, err = p.Cfg.Query.InsertStockPrice(stockData, context.TODO())
		if err != nil {
			fmt.Printf("Error insertting stock price for symbol: %s", symbol.Symbol)
			continue
		}
	}
	return nil
}

func HandlerAddNewSymbol(p *Command) error {
	if len(p.Input) != 1 {
		return errors.New("Please provide a stonk symbol")
	}
	symbol := p.Input[0]
	stocks := []models.StockData{}
	stockChan := make(chan models.StockData)
	errChan := make(chan error)
	dates := []time.Time{}
	daysBefore := 1

	for len(dates) != p.Cfg.HistoricalTimeFrame {
		prevDate := time.Now().Add(-time.Hour * 24 * time.Duration(daysBefore))

		if prevDate.Weekday() != time.Sunday && prevDate.Weekday() != time.Saturday {
			dates = append(dates, prevDate)
		}
		daysBefore++
	}

	for i := 0; i < p.Cfg.HistoricalTimeFrame; i++ {
		go func() {
			stockData, err := p.Cfg.ApiClient.GetPrices(symbol, dates[i].Format("2006-01-02"))
			if err != nil { // holiday will cause error, so it is ok to swallow
				errChan <- err
				return
			}
			stockChan <- stockData
		}()
	}
	count := 0
	ended := false
	for !ended {
		select {
		case err := <-errChan:
			count++
			fmt.Printf("Error fetching stock price for symbol for %d days: %s\n %v\n", p.Cfg.HistoricalTimeFrame, symbol, err)
			if count == p.Cfg.HistoricalTimeFrame {
				ended = true
			}
		case s := <-stockChan:
			count++
			stocks = append(stocks, s)
			if count == p.Cfg.HistoricalTimeFrame {
				ended = true
			}
		}
	}
	_, err := p.Cfg.Query.InsertStockPrices(stocks, symbol, context.TODO())
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
