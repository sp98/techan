package techan

import "github.com/sdcoffey/big"

type bearishHaramiIndicator struct {
	series *TimeSeries
}

// NewBearishHaramiIndicator returns an indicator that returns big.One if bearish harami pattern is formed, else
// returns big.ZERO. http://www.investopedia.com/terms/b/bearishharami.asp
func NewBearishHaramiIndicator(series *TimeSeries) Indicator {
	return bearishHaramiIndicator{series}
}

func (b bearishHaramiIndicator) Calculate(index int) big.Decimal {
	if index < 1 {
		// Harami is a 2-candle pattern
		return big.ZERO
	}

	if NewBullishIndicator(b.series).Calculate(index-1) == big.ONE && NewBearishIndicator(b.series).Calculate(index) == big.ONE {
		prevOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index - 1)
		prevClosePrice := NewClosePriceIndicator(b.series).Calculate(index - 1)
		currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
		currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)

		if currOpenPrice.GT(prevOpenPrice) && currOpenPrice.LT(prevClosePrice) && currClosePrice.GT(prevOpenPrice) && currClosePrice.LT(prevClosePrice) {
			return big.ONE
		}
	}

	return big.ZERO
}
