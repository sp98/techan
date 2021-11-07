package techan

import "github.com/sdcoffey/big"

type realBodyIndicator struct {
	series *TimeSeries
}

func NewRealBodyIndicator(series *TimeSeries) Indicator {
	return realBodyIndicator{series}
}

func (r realBodyIndicator) Calculate(index int) big.Decimal {
	candle := r.series.Candles[index]
	return (candle.ClosePrice.Sub(candle.OpenPrice)).Abs()
}
