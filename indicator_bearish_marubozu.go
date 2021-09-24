package techan

import "github.com/sdcoffey/big"

type bearishMarubozuIndicator struct {
	series *TimeSeries
}

// NewBearishMarubozuIndicator returns and indicator with returns big.ONE if candle is bearish marubozu, else returns big.ZERO
func NewBearishMarubozuIndicator(series *TimeSeries) Indicator {
	return bearishMarubozuIndicator{series}
}

func (b bearishMarubozuIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
	currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)
	currHighPrice := NewHighPriceIndicator(b.series).Calculate(index)
	currLowPrice := NewLowPriceIndicator(b.series).Calculate(index)

	// TODO: should we compare the candle height with average of pervious candles?
	if currClosePrice.LT(currOpenPrice) && currOpenPrice.EQ(currHighPrice) && currClosePrice.EQ(currLowPrice) {
		return big.ONE
	}

	return big.ZERO
}
