package techan

import "github.com/sdcoffey/big"

type adxIndicator struct {
	series         *TimeSeries
	avgDXIndicator Indicator
}

func NewADXIndicator(series *TimeSeries, diWindow, adxWndow int) Indicator {
	return adxIndicator{
		series:         series,
		avgDXIndicator: NewMMAIndicator(NewDXIndicator(series, diWindow), adxWndow),
	}
}

func (a adxIndicator) Calculate(index int) big.Decimal {
	return a.avgDXIndicator.Calculate(index)
}
