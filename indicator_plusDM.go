package techan

import "github.com/sdcoffey/big"

type plusDMIndicator struct {
	series *TimeSeries
}

func NewPlusDMIndicator(series *TimeSeries) Indicator {
	return plusDMIndicator{series}
}

func (p plusDMIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	currentCandle := p.series.Candles[index]
	prevCandle := p.series.Candles[index-1]

	upMove := currentCandle.MaxPrice.Sub(prevCandle.MaxPrice)
	downMove := prevCandle.MinPrice.Sub(currentCandle.MinPrice)

	if upMove.GT(downMove) && upMove.GT(big.NewFromInt(0)) {
		return upMove
	}

	return big.ZERO
}
