package calculation

// https://www.litefinance.org/blog/for-beginners/types-of-forex-charts/heikin-ashi-candles/
import (
	"math"
)

type PriceCal struct {
	Open  float64
	Close float64
	High  float64
	Low   float64
}

func GetHeikinDailyOpen(prev *PriceCal) float64 {
	return (prev.Open + prev.Close) / 2
}

func GetHeikinDailyClose(price *PriceCal) float64 {
	return (price.Open + price.Close + price.High + price.Low) / 4
}

func GetHeikinDailyHigh(price *PriceCal, prev *PriceCal) float64 {
	open := GetHeikinDailyOpen(prev)
	close := GetHeikinDailyClose(price)
	return math.Max(math.Max(open, close), price.High)
}

func GetHeikinDailyLow(price *PriceCal, prev *PriceCal) float64 {
	open := GetHeikinDailyOpen(prev)
	close := GetHeikinDailyClose(price)
	return math.Min(math.Min(open, close), price.Low)
}

// prev should be heikin, price should be normal
func GetIsSpinningTop(price *PriceCal, prev *PriceCal) bool {
	open := GetHeikinDailyOpen(prev)
	high := GetHeikinDailyHigh(price, prev)
	low := GetHeikinDailyLow(price, prev)
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
	high := GetHeikinDailyHigh(price, prev)
	low := GetHeikinDailyLow(price, prev)
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
	high := GetHeikinDailyHigh(price, prev)
	low := GetHeikinDailyLow(price, prev)
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

func GetIsBull(price *PriceCal, prev *PriceCal) bool {
	open := GetHeikinDailyOpen(prev)
	low := GetHeikinDailyLow(price, prev)
	close := GetHeikinDailyClose(price)

	return open < close && int(open*100) == int(low*100)
}

func GetIsBear(price *PriceCal, prev *PriceCal) bool {
	open := GetHeikinDailyOpen(prev)
	high := GetHeikinDailyHigh(price, prev)
	close := GetHeikinDailyClose(price)

	return close < open && int(open*100) == int(high*100)
}
