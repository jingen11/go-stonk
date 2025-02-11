package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/jingen11/stonk-tracker/internal/utils"
)

type Params struct {
	cfg   utils.ProjectConfig
	input []string
}

func HandleRefresh(p *Params) error {
	symbols, err := p.cfg.Query.GetAllSymbols(context.TODO())
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		stockData, err := p.cfg.ApiClient.GetPrices(symbol.Symbol)
		if err != nil {
			fmt.Printf("Error getting stockData for symbol: %s", symbol.Symbol)
			continue
		}
		_, err = p.cfg.Query.InsertStockPrices(stockData, context.TODO())
		if err != nil {
			fmt.Printf("Error insertting stock price for symbol: %s", symbol.Symbol)
			continue
		}
	}
	return nil
}

func HandlerAddNewSymbol(p *Params) error {
	if len(p.input) != 1 {
		return errors.New("Please provide a stonk symbol")
	}
	symbol := p.input[0]
	stockData, err := p.cfg.ApiClient.GetPrices(symbol)
	if err != nil {
		fmt.Printf("Error getting stockData for symbol: %s", symbol)
		return err

	}
	_, err = p.cfg.Query.InsertStockPrices(stockData, context.TODO())
	if err != nil {
		fmt.Printf("Error insertting stock price for symbol: %s", symbol)
		return err
	}
	return nil
}

// func HandleGetTrendReversal(p *Params) error {
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
