package techan

import "github.com/sdcoffey/big"

type bearishMarubozuIndicator struct {
	series *TimeSeries
	offset big.Decimal
}

// NewBearishMarubozuIndicator returns and indicator with returns big.ONE if candle is bearish marubozu, else returns big.ZERO
func NewBearishMarubozuIndicator(series *TimeSeries, offset big.Decimal) Indicator {
	return bearishMarubozuIndicator{series, offset}
}

func (b bearishMarubozuIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	// Return if candle is not bearish
	if NewBearishIndicator(b.series).Calculate(index) != big.ONE {
		return big.ZERO
	}

	currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
	currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)
	currHighPrice := NewHighPriceIndicator(b.series).Calculate(index)
	currLowPrice := NewLowPriceIndicator(b.series).Calculate(index)

	isPureBearishMarubozu := currOpenPrice.EQ(currHighPrice) && currClosePrice.EQ(currLowPrice)
	isBearishMarubuzo := (((currOpenPrice.Sub(currClosePrice)).Div(currHighPrice.Sub(currLowPrice))).Mul(big.NewFromInt(100))).GTE(b.offset)

	// TODO: should we compare the candle height with average of pervious candles?
	if isPureBearishMarubozu || isBearishMarubuzo {
		return big.ONE
	}

	return big.ZERO
}
