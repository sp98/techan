package techan

import (
	"github.com/sdcoffey/big"
)

type parabolicSARIndicator struct {
	highPriceIndicator    Indicator
	lowPriceIndicator     Indicator
	closePriceIndicator   Indicator
	maxAcceleration       big.Decimal
	accelerationFactor    big.Decimal
	accelerationIncrement big.Decimal
	accelerationStart     big.Decimal
	currentExtremePoint   big.Decimal
	minMaxExtremePoint    big.Decimal
	sar                   big.Decimal
	window                int
	resultCache           resultCache
	currentTrend          bool // uptrend if true, downstrend if false
	startTrendIndex       int
}

func NewParabolicSARIndicator(series *TimeSeries, aF, maxA, increment big.Decimal) Indicator {
	return &parabolicSARIndicator{
		highPriceIndicator:    NewHighPriceIndicator(series),
		lowPriceIndicator:     NewLowPriceIndicator(series),
		closePriceIndicator:   NewClosePriceIndicator(series),
		accelerationFactor:    aF,
		maxAcceleration:       maxA,
		accelerationIncrement: increment,
		accelerationStart:     aF,
		window:                20,
		resultCache:           make([]*big.Decimal, 10000),
	}
}

func (p *parabolicSARIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(p, index, func(i int) big.Decimal {
		p.sar = p.lowPriceIndicator.Calculate(index)
		p.currentExtremePoint = p.highPriceIndicator.Calculate(index)
		return p.sar
	}); cachedValue != nil {
		return *cachedValue
	}

	if index == 1 {
		p.currentTrend = p.closePriceIndicator.Calculate(index - 1).LT(p.closePriceIndicator.Calculate(index))
		if !p.currentTrend {
			p.sar = NewHighestValueIndicator(p.highPriceIndicator, 2).Calculate(index)
			p.currentExtremePoint = p.sar
			p.minMaxExtremePoint = p.currentExtremePoint
		} else {
			p.sar = NewLowestValueIndicator(p.lowPriceIndicator, 2).Calculate(index)
			p.currentExtremePoint = p.sar
			p.minMaxExtremePoint = p.currentExtremePoint
		}

		cacheResult(p, index, p.sar)
		return p.sar
	}

	prevSar := p.Calculate(index - 1)
	if p.currentTrend { // if up trend
		p.sar = prevSar.Add(p.accelerationFactor.Mul((p.currentExtremePoint.Sub(prevSar))))
		p.currentTrend = p.lowPriceIndicator.Calculate(index).GT(p.sar)
		if !p.currentTrend {
			if p.minMaxExtremePoint.GT(p.highPriceIndicator.Calculate(index)) {
				p.sar = p.minMaxExtremePoint
			} else {
				p.sar = p.highPriceIndicator.Calculate(index)
			}
			p.currentTrend = false
			p.startTrendIndex = index
			p.accelerationFactor = p.accelerationStart
			p.currentExtremePoint = p.lowPriceIndicator.Calculate(index)
			p.minMaxExtremePoint = p.currentExtremePoint
		} else {
			lowestPriceOfTwoPreviousBars := NewLowestValueIndicator(p.lowPriceIndicator, minInt(2, (index-p.startTrendIndex))).Calculate(index - 1)
			if p.sar.GT(lowestPriceOfTwoPreviousBars) {
				p.sar = lowestPriceOfTwoPreviousBars
			}

			p.currentExtremePoint = NewHighestValueIndicator(p.highPriceIndicator, (index - p.startTrendIndex + 1)).Calculate(index)
			if p.currentExtremePoint.GT(p.minMaxExtremePoint) {
				p.incrementAcceleration()
				p.minMaxExtremePoint = p.currentExtremePoint
			}
		}
	} else {
		p.sar = prevSar.Sub(p.accelerationFactor.Mul((prevSar.Sub(p.currentExtremePoint))))
		p.currentTrend = p.highPriceIndicator.Calculate(index).GTE(p.sar)
		if p.currentTrend {
			if p.minMaxExtremePoint.LT(p.lowPriceIndicator.Calculate(index)) {
				p.sar = p.minMaxExtremePoint
			} else {
				p.sar = p.lowPriceIndicator.Calculate(index)
			}
			p.accelerationFactor = p.accelerationStart
			p.startTrendIndex = index
			p.currentExtremePoint = p.highPriceIndicator.Calculate(index)
			p.minMaxExtremePoint = p.currentExtremePoint
		} else {
			highestPriceOfTwoPreviousBars := NewHighestValueIndicator(p.highPriceIndicator, minInt(2, index-p.startTrendIndex)).Calculate(index - 1)
			if p.sar.LT(highestPriceOfTwoPreviousBars) {
				p.sar = highestPriceOfTwoPreviousBars
			}
			p.currentExtremePoint = NewLowestValueIndicator(p.lowPriceIndicator, (index - p.startTrendIndex + 1)).Calculate(index)
			if p.currentExtremePoint.LT(p.minMaxExtremePoint) {
				p.incrementAcceleration()
				p.minMaxExtremePoint = p.currentExtremePoint
			}
		}
	}

	cacheResult(p, index, p.sar)
	return p.sar
}

func (p *parabolicSARIndicator) incrementAcceleration() {
	if p.accelerationFactor.GTE(p.maxAcceleration) {
		p.accelerationFactor = p.maxAcceleration
	} else {
		p.accelerationFactor = p.accelerationFactor.Add(p.accelerationIncrement)
	}
}

func (p parabolicSARIndicator) cache() resultCache {
	return p.resultCache
}

func (p parabolicSARIndicator) setCache(cache resultCache) {
	p.resultCache = cache
}

func (p parabolicSARIndicator) windowSize() int {
	return p.window
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
