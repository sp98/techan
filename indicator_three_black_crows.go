package techan

import "github.com/sdcoffey/big"

type threeBlackCrows struct {
	series                  *TimeSeries
	lowerShadowIndicator    Indicator
	avgLowerShadowIndicator Indicator
	factor                  big.Decimal
	bullishCandleIndex      int
}

func NewThreeBlackCrowsIndicator(series *TimeSeries, window int, factor big.Decimal) Indicator {
	return threeBlackCrows{
		series:                  series,
		lowerShadowIndicator:    NewLowerShadowIndicator(series),
		avgLowerShadowIndicator: NewSimpleMovingAverage(NewLowerShadowIndicator(series), window),
		factor:                  factor,
	}
}

func (t threeBlackCrows) Calculate(index int) big.Decimal {
	if index < 3 {
		// We need 4 candles: 1 black, 3 white
		return big.ZERO
	}
	t.bullishCandleIndex = index - 3
	isBullish := NewBullishIndicator(t.series).Calculate(t.bullishCandleIndex).EQ(big.NewFromInt(1))

	if isBullish && t.isBlackCrow(index-2) && t.isBlackCrow(index-1) && t.isBlackCrow(index) {
		return big.ONE
	}
	return big.ZERO
}

func (t threeBlackCrows) isBlackCrow(index int) bool {
	if NewBearishIndicator(t.series).Calculate(index).EQ(big.NewFromInt(1)) {
		if NewBullishIndicator(t.series).Calculate(index - 1).EQ(big.NewFromInt(1)) {
			// first soldier
			currCandle := t.series.Candles[index]
			prevCandle := t.series.Candles[index-1]
			return t.hasVeryShortLowerShadow(index) && currCandle.OpenPrice.LT(prevCandle.MaxPrice)
		} else {
			return t.hasVeryShortLowerShadow(index) && t.isDeclining(index)
		}
	}

	return false
}

func (t threeBlackCrows) isDeclining(index int) bool {
	currCandle := t.series.Candles[index]
	prevCandle := t.series.Candles[index-1]
	currOpenPrice := currCandle.OpenPrice
	currClosePrice := currCandle.ClosePrice
	prevOpenPrice := prevCandle.OpenPrice
	prevClosePrice := prevCandle.ClosePrice

	return currOpenPrice.LT(prevOpenPrice) && currOpenPrice.GT(prevClosePrice) && currClosePrice.LT(prevClosePrice)
}

func (t threeBlackCrows) hasVeryShortLowerShadow(index int) bool {
	currLowerShadow := t.lowerShadowIndicator.Calculate(index)
	avgLowerShadow := t.avgLowerShadowIndicator.Calculate(t.bullishCandleIndex)
	return currLowerShadow.LT(avgLowerShadow.Mul(t.factor))
}
