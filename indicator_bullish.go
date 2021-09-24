package techan

import "github.com/sdcoffey/big"

type bullishIndicator struct {
	series *TimeSeries
}

// NewBullishIndicator returns an indicator which returns big.ONE if candle is bullish else returns big.ZERO
func NewBullishIndicator(series *TimeSeries) Indicator {
	return bullishIndicator{series}
}

func (b bullishIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	openPrice := NewOpenPriceIndicator(b.series).Calculate(index)
	closePrice := NewClosePriceIndicator(b.series).Calculate(index)

	if openPrice.LT(closePrice) {
		return big.ONE
	}

	return big.ZERO
}
