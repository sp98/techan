package techan

import "github.com/sdcoffey/big"

type bullishEngulfingIndicator struct {
	series *TimeSeries
}

func NewBullishEngulfingIndicator(series *TimeSeries) Indicator {
	return bullishEngulfingIndicator{series}
}

func (b bullishEngulfingIndicator) Calculate(index int) big.Decimal {
	if index < 1 {
		// Engulfing is a 2-candle pattern
		return big.ZERO
	}

	if NewBearishIndicator(b.series).Calculate(index-1) == big.ONE && NewBullishIndicator(b.series).Calculate(index) == big.ONE {
		prevOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index - 1)
		prevClosePrice := NewClosePriceIndicator(b.series).Calculate(index - 1)
		currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
		currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)

		if currOpenPrice.LT(prevOpenPrice) && currOpenPrice.LT(prevClosePrice) && currClosePrice.GT(prevOpenPrice) && currClosePrice.GT(prevClosePrice) {
			return big.ONE
		}
	}

	return big.ZERO
}
