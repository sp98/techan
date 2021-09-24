package techan

import "github.com/sdcoffey/big"

type bullishMarubozuIndicator struct {
	series *TimeSeries
}

// NewBullishMarubozuIndicator returns and indicator with returns big.ONE if candle is bullish marubozu, else returns big.ZERO
func NewBullishMarubozuIndicator(series *TimeSeries) Indicator {
	return bullishMarubozuIndicator{series}
}

func (b bullishMarubozuIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
	currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)
	currHighPrice := NewHighPriceIndicator(b.series).Calculate(index)
	currLowPrice := NewLowPriceIndicator(b.series).Calculate(index)

	// TODO: should we compare the candle height with average of pervious candles?
	if currOpenPrice.LT(currClosePrice) && currOpenPrice.EQ(currLowPrice) && currClosePrice.EQ(currHighPrice) {
		return big.ONE
	}

	return big.ZERO
}
