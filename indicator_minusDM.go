package techan

import "github.com/sdcoffey/big"

type minusDMIndicator struct {
	series *TimeSeries
}

func NewMinusDMIndicator(series *TimeSeries) Indicator {
	return minusDMIndicator{series}
}

func (m minusDMIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	currentCandle := m.series.Candles[index]
	prevCandle := m.series.Candles[index-1]

	upMove := currentCandle.MaxPrice.Sub(prevCandle.MaxPrice)
	downMove := prevCandle.MinPrice.Sub(currentCandle.MinPrice)

	if downMove.GT(upMove) && downMove.GT(big.NewFromInt(0)) {
		return downMove
	}

	return big.ZERO
}
