package techan

import "github.com/sdcoffey/big"

type bullishHaramiIndicator struct {
	series *TimeSeries
}

// NewBullishHaramiIndicator returns an indicator that returns big.One if bullish harami pattern is formed, else
// returns big.ZERO. http://www.investopedia.com/terms/b/bullishharami.asp
func NewBullishHaramiIndicator(series *TimeSeries) Indicator {
	return bullishHaramiIndicator{series}
}

func (b bullishHaramiIndicator) Calculate(index int) big.Decimal {
	if index < 1 {
		// Harami is a 2-candle pattern
		return big.ZERO
	}

	if NewBearishIndicator(b.series).Calculate(index-1) == big.ONE && NewBullishIndicator(b.series).Calculate(index) == big.ONE {
		prevOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index - 1)
		prevClosePrice := NewClosePriceIndicator(b.series).Calculate(index - 1)
		currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
		currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)

		if currOpenPrice.LT(prevOpenPrice) && currOpenPrice.GT(prevClosePrice) && currClosePrice.LT(prevOpenPrice) && currClosePrice.GT(prevClosePrice) {
			return big.ONE
		}
	}

	return big.ZERO
}
