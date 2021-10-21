package techan

import "github.com/sdcoffey/big"

type minusDIIndicator struct {
	avgMinusDMIndicator Indicator
	atrIndicator        Indicator
	window              int
}

func NewMinusDIIndicator(series *TimeSeries, window int) Indicator {
	return minusDIIndicator{
		avgMinusDMIndicator: NewMMAIndicator(NewMinusDMIndicator(series), window),
		atrIndicator:        NewAverageTrueRangeIndicator(series, window),
		window:              window,
	}
}

func (m minusDIIndicator) Calculate(index int) big.Decimal {
	return m.avgMinusDMIndicator.Calculate(index).Div(m.atrIndicator.Calculate(index)).Mul(big.NewFromInt(100))
}
