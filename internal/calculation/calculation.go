package calculation

import (
	"math"
	"strconv"

	"github.com/jingen11/stonk-tracker/internal/models"
)

type PriceCal struct {
	Open  float64
	Close float64
	High  float64
	Low   float64
}

func ConvertPrice(price *models.Price) (PriceCal, error) {
	priceCal := PriceCal{}
	open, err := strconv.ParseFloat(price.Open, 64)
	if err != nil {
		return priceCal, err
	}
	close, err := strconv.ParseFloat(price.Close, 64)
	if err != nil {
		return priceCal, err
	}
	high, err := strconv.ParseFloat(price.High, 64)
	if err != nil {
		return priceCal, err
	}
	low, err := strconv.ParseFloat(price.Low, 64)
	if err != nil {
		return priceCal, err
	}
	priceCal.Open = open
	priceCal.Close = close
	priceCal.High = high
	priceCal.Low = low
	return priceCal, nil
}

func GetHeikinDailyClose(price *PriceCal) float64 {
	return (price.Open + price.Close + price.High + price.Low) / 4
}

func GetHeikinDailyHigh(price *PriceCal) float64 {
	return math.Max(math.Max(price.Open, price.Close), price.High)
}

func GetHeikinDailyLow(price *PriceCal) float64 {
	return math.Min(math.Max(price.Open, price.Close), price.Low)
}

func GetHeikinDailyOpen(prev *PriceCal) float64 {
	return (prev.Open + prev.Close) / 2
}

// prev should be heikin, price should be normal
func GetIsSpinningTop(price *PriceCal, prev *PriceCal) bool {
	open := GetHeikinDailyOpen(prev)
	high := GetHeikinDailyHigh(price)
	low := GetHeikinDailyLow(price)
	close := GetHeikinDailyClose(price)

	dailyRange := high - low

	diff := close - open

	candleRange := math.Abs(diff)

	percentage := candleRange / dailyRange * 100

	if percentage < 20 && percentage > 5 {
		if int(high*100) == int(close*100) || int(high*100) == int(open*100) {
			return false
		}

		if int(low*100) == int(close*100) || int(low*100) == int(open*100) {
			return false
		}

		return true
	}

	return false
}

func GetIsDojiStar(price *PriceCal, prev *PriceCal) bool {
	open := GetHeikinDailyOpen(prev)
	high := GetHeikinDailyHigh(price)
	low := GetHeikinDailyLow(price)
	close := GetHeikinDailyClose(price)

	dailyRange := high - low

	diff := close - open

	candleRange := math.Abs(diff)

	percentage := candleRange / dailyRange * 100

	if percentage < 5 {
		return true //https://www.investopedia.com/terms/d/doji.asp
	}

	return false
}

func GetIsGravestoneDoji(price *PriceCal, prev *PriceCal) bool {
	high := GetHeikinDailyHigh(price)
	low := GetHeikinDailyLow(price)
	close := GetHeikinDailyClose(price)
	open := GetHeikinDailyOpen(prev)

	dailyRange := high - low

	diff := close - open

	candleRange := math.Abs(diff)

	percentage := candleRange / dailyRange * 100

	if percentage < 5 {
		gravestoneRange := math.Abs(close - low)
		gPercentage := gravestoneRange / dailyRange * 100
		if gPercentage < 5 {
			return true
		}
		return false
	}

	return false
}

func GetIsUptrend(price *PriceCal, prev *PriceCal) bool {
	close := GetHeikinDailyClose(price)
	open := GetHeikinDailyOpen(prev)

	diff := close - open

	uptrend := true

	if diff < 0 {
		uptrend = false
	}

	return uptrend
}
