package techan

import "github.com/sdcoffey/big"

type UpperShadowIndicator struct {
	series *TimeSeries
}

func NewUpperShadowIndicator(series *TimeSeries) Indicator {
	return UpperShadowIndicator{series}
}

func (u UpperShadowIndicator) Calculate(index int) big.Decimal {
	candle := u.series.Candles[index]
	openPrice := candle.OpenPrice
	closePrice := candle.ClosePrice
	highPrice := candle.MaxPrice

	if closePrice.GT(openPrice) {
		// Bullish
		return highPrice.Sub(closePrice)

	}

	// Bearish
	return highPrice.Sub(openPrice)
}
