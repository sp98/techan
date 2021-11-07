package techan

import "github.com/sdcoffey/big"

type threeWhiteSoldiers struct {
	series                  *TimeSeries
	upperShadowIndicator    Indicator
	avgUpperShadowIndicator Indicator
	factor                  big.Decimal
	bearishCandleIndex      int
}

func NewThreeWhiteSoldiersIndicator(series *TimeSeries, window int, factor big.Decimal) Indicator {
	return threeWhiteSoldiers{
		series:                  series,
		upperShadowIndicator:    NewUpperShadowIndicator(series),
		avgUpperShadowIndicator: NewSimpleMovingAverage(NewUpperShadowIndicator(series), window),
		factor:                  factor,
	}
}

func (t threeWhiteSoldiers) Calculate(index int) big.Decimal {
	if index < 3 {
		// We need 4 candles: 1 black, 3 white
		return big.ZERO
	}
	t.bearishCandleIndex = index - 3
	isBearish := NewBearishIndicator(t.series).Calculate(t.bearishCandleIndex).EQ(big.NewFromInt(1))

	if isBearish && t.isWhiteSoldier(index-2) && t.isWhiteSoldier(index-1) && t.isWhiteSoldier(index) {
		return big.ONE
	}
	return big.ZERO
}

func (t threeWhiteSoldiers) isWhiteSoldier(index int) bool {
	if NewBullishIndicator(t.series).Calculate(index).EQ(big.NewFromInt(1)) {
		if NewBearishIndicator(t.series).Calculate(index - 1).EQ(big.NewFromInt(1)) {
			// first soldier
			currCandle := t.series.Candles[index]
			prevCandle := t.series.Candles[index-1]
			return t.hasVeryShortUpperShadow(index) && currCandle.OpenPrice.GT(prevCandle.MinPrice)
		} else {
			return t.hasVeryShortUpperShadow(index) && t.isGrowing(index)
		}
	}

	return false
}

func (t threeWhiteSoldiers) isGrowing(index int) bool {
	currCandle := t.series.Candles[index]
	prevCandle := t.series.Candles[index-1]
	currOpenPrice := currCandle.OpenPrice
	currClosePrice := currCandle.ClosePrice
	prevOpenPrice := prevCandle.OpenPrice
	prevClosePrice := prevCandle.ClosePrice

	return currOpenPrice.GT(prevOpenPrice) && currOpenPrice.LT(prevClosePrice) && currClosePrice.GT(prevClosePrice)
}

func (t threeWhiteSoldiers) hasVeryShortUpperShadow(index int) bool {
	currUpperShadow := t.upperShadowIndicator.Calculate(index)
	avgUpperShadow := t.avgUpperShadowIndicator.Calculate(t.bearishCandleIndex)
	return currUpperShadow.LT(avgUpperShadow.Mul(t.factor))
}
