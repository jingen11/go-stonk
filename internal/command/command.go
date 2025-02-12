package command

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jingen11/stonk-tracker/internal/calculation"
	"github.com/jingen11/stonk-tracker/internal/db"
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

func HandleGetInfo(p *Command) error {
	symbols, err := p.Cfg.Query.GetAllSymbols(context.TODO())
	if err != nil {
		return err
	}
	limit := 80
	endChan := make(chan bool)
	for _, s := range symbols {
		go getSymbolInfo(p, s, limit, endChan)
	}

	for range symbols {
		<-endChan
	}
	return nil
}

func getSymbolInfo(p *Command, s models.Symbol, limit int, endChan chan bool) {
	prices, err := p.Cfg.Query.GetStockPrices(context.TODO(), &db.GetStockPriceOpt{
		Symbol: s.Symbol,
		Limit:  int64(limit),
	})

	if err != nil {
		fmt.Printf("Cannot get info for symbol: %s\n", s.Symbol)
		endChan <- false
		return
	}

	if len(prices) != limit {
		fmt.Printf("Insufficient data point for symbol: %s\n", s.Symbol)
		endChan <- false
		return
	}

	slices.Reverse(prices)

	pr := models.Price{}

	for i, p := range prices {
		if i == 0 {
			pr = p
		}
		price := calculation.PriceCal{
			Open:  p.Open,
			Close: p.Close,
			High:  p.High,
			Low:   p.Low,
		}

		prev := calculation.PriceCal{
			Open:  pr.Open,  // HA
			Close: pr.Close, // HA
		}
		po := calculation.GetHeikinDailyOpen(&prev)
		pc := calculation.GetHeikinDailyClose(&price)

		pr = models.Price{
			Open:  po,
			Close: pc,
		}
	}

	price := calculation.PriceCal{
		Open:  prices[limit-1].Open,
		Close: prices[limit-1].Close,
		High:  prices[limit-1].High,
		Low:   prices[limit-1].Low,
	}

	prev := calculation.PriceCal{
		Open:  pr.Open,  // HA
		Close: pr.Close, // HA
	}

	o := calculation.GetHeikinDailyOpen(&prev)
	c := calculation.GetHeikinDailyClose(&price)
	h := calculation.GetHeikinDailyHigh(&price, &prev)
	l := calculation.GetHeikinDailyLow(&price, &prev)
	u := calculation.GetIsUptrend(&price, &prev)
	bu := calculation.GetIsBull(&price, &prev)
	be := calculation.GetIsBear(&price, &prev)
	st := calculation.GetIsSpinningTop(&price, &prev)
	ds := calculation.GetIsDojiStar(&price, &prev)
	g := calculation.GetIsGravestoneDoji(&price, &prev)
	sen := getSentiment(u, bu, be, st, ds, g)

	fmt.Printf("------------------------------------\nDate: %s\nSymbol: %s\nOHLC: %.2f, %.2f, %.2f, %.2f\nUptrend: %v\nBull: %v\nBear: %v\nSpinningTop: %v\nDoji: %v\nGrave: %v \nSentiment: %s\n",
		prices[limit-1].Date.Time().Format("2006-01-02"), s.Symbol, o, h, l, c, u, bu, be, st, ds, g, sen)
	endChan <- true
}

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

func getSentiment(u, bu, be, st, ds, g bool) string {
	if !u && ds {
		return "buy"
	}
	if u && ds {
		return "sell"
	}
	if u && st {
		return "hold, sell"
	}
	if be {
		return "SELL"
	}
	if bu {
		return "hold, add"
	}
	if g {
		return "sell"
	}
	if u && st {
		return "hold"
	}
	if !u {
		return "sell"
	}
	if u {
		return "hold"
	}
	return "no action"
}
