package techan

import "github.com/sdcoffey/big"

type lowerShadowIndicator struct {
	series *TimeSeries
}

func NewLowerShadowIndicator(series *TimeSeries) Indicator {
	return lowerShadowIndicator{series}
}

func (l lowerShadowIndicator) Calculate(index int) big.Decimal {
	candle := l.series.Candles[index]
	openPrice := candle.OpenPrice
	closePrice := candle.ClosePrice
	lowPrice := candle.MinPrice

	if closePrice.GT(openPrice) {
		// Bullish
		return openPrice.Sub(lowPrice)

	}

	// Bearish
	return closePrice.Sub(lowPrice)
}
