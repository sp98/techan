package techan

import "github.com/sdcoffey/big"

type dojiIndicator struct {
	bodyHeightIndicator    Indicator
	avgBodyHeightIndicator Indicator
	factor                 big.Decimal
}

func NewDojiIndicator(series *TimeSeries, window int, bodyFactor big.Decimal) Indicator {
	return dojiIndicator{
		bodyHeightIndicator:    NewRealBodyIndicator(series),
		avgBodyHeightIndicator: NewSimpleMovingAverage(NewRealBodyIndicator(series), window),
		factor:                 bodyFactor,
	}
}

func (d dojiIndicator) Calculate(index int) big.Decimal {
	currentCandleBodyHeight := d.bodyHeightIndicator.Calculate(index)
	avgBodyHeight := d.avgBodyHeightIndicator.Calculate(index - 1)

	if currentCandleBodyHeight.LT(avgBodyHeight.Mul(d.factor)) {
		return big.ONE
	}

	return big.ZERO
}
