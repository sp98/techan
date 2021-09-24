package techan

import "github.com/sdcoffey/big"

type bearishIndicator struct {
	series *TimeSeries
}

// bearishIndicator returns an indicator which returns big.ONE if candle is bullish else returns big.ZERO
func NewBearishIndicator(series *TimeSeries) Indicator {
	return bearishIndicator{series}
}

func (b bearishIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	openPrice := NewOpenPriceIndicator(b.series).Calculate(index)
	closePrice := NewClosePriceIndicator(b.series).Calculate(index)

	if openPrice.GT(closePrice) {
		return big.ONE
	}

	return big.ZERO
}
