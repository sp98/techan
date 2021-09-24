package techan

import "github.com/sdcoffey/big"

type bearishEngulfingIndicator struct {
	series *TimeSeries
}

func NewBearishEngulfingIndicator(series *TimeSeries) Indicator {
	return bearishEngulfingIndicator{series}
}

func (b bearishEngulfingIndicator) Calculate(index int) big.Decimal {
	if index < 1 {
		// Engulfing is a 2-candle pattern
		return big.ZERO
	}

	if NewBearishIndicator(b.series).Calculate(index) == big.ONE && NewBullishIndicator(b.series).Calculate(index-1) == big.ONE {
		prevOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index - 1)
		prevClosePrice := NewClosePriceIndicator(b.series).Calculate(index - 1)
		currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
		currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)

		if currOpenPrice.GT(prevOpenPrice) && currOpenPrice.GT(prevClosePrice) && currClosePrice.LT(prevOpenPrice) && currClosePrice.LT(prevClosePrice) {
			return big.ONE
		}
	}

	return big.ZERO
}
