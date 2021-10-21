package techan

import "github.com/sdcoffey/big"

type plusDIIndicator struct {
	avgPlusDMIndicator Indicator
	atrIndicator       Indicator
	window             int
}

func NewPlusDIIndicator(series *TimeSeries, window int) Indicator {
	return plusDIIndicator{
		avgPlusDMIndicator: NewMMAIndicator(NewPlusDMIndicator(series), window),
		atrIndicator:       NewAverageTrueRangeIndicator(series, window),
		window:             window,
	}
}

func (p plusDIIndicator) Calculate(index int) big.Decimal {
	return p.avgPlusDMIndicator.Calculate(index).Div(p.atrIndicator.Calculate(index)).Mul(big.NewFromInt(100))
}
