package techan

import (
	"github.com/sdcoffey/big"
)

type bullishMarubozuIndicator struct {
	series *TimeSeries
	offset big.Decimal
}

// NewBullishMarubozuIndicator returns and indicator with returns big.ONE if candle is bullish marubozu, else returns big.ZERO
func NewBullishMarubozuIndicator(series *TimeSeries, offset big.Decimal) Indicator {
	return bullishMarubozuIndicator{series, offset}
}

func (b bullishMarubozuIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	// Return if candle is not bullish
	if NewBullishIndicator(b.series).Calculate(index) != big.ONE {
		return big.ZERO
	}

	currOpenPrice := NewOpenPriceIndicator(b.series).Calculate(index)
	currClosePrice := NewClosePriceIndicator(b.series).Calculate(index)
	currHighPrice := NewHighPriceIndicator(b.series).Calculate(index)
	currLowPrice := NewLowPriceIndicator(b.series).Calculate(index)

	// TODO: should we compare the candle height with average of pervious candles?

	isPureBullishMarubozu := currOpenPrice.EQ(currLowPrice) && currClosePrice.EQ(currHighPrice)
	isBullishMarubozu := ((currClosePrice.Sub(currOpenPrice).Div(currHighPrice.Sub(currLowPrice))).Mul(big.NewFromInt(100))).GTE(b.offset)
	if isPureBullishMarubozu || isBullishMarubozu {
		return big.ONE
	}

	return big.ZERO
}
