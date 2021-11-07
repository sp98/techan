package techan

import (
	"github.com/sdcoffey/big"
)

type hiekinAshiSmoothedRSIIndicator struct {
	rsiIndicator Indicator
	window       int
	resultCache  resultCache
	isSmooth     bool
}

func NewHiekinAshiSmoothedRSIIndicator(indicator Indicator, timeframe int, isSmooth bool) Indicator {
	return hiekinAshiSmoothedRSIIndicator{
		rsiIndicator: NewRelativeStrengthIndexIndicator(indicator, timeframe),
		resultCache:  make([]*big.Decimal, 10000),
		window:       timeframe,
		isSmooth:     isSmooth,
	}
}

func (hasr hiekinAshiSmoothedRSIIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(hasr, index, func(i int) big.Decimal {
		return hasr.rsiIndicator.Calculate(index).Sub(big.NewFromInt(50))
	}); cachedValue != nil {
		return *cachedValue
	}

	todayVal := hasr.rsiIndicator.Calculate(index).Sub(big.NewFromInt(50))
	lastVal := hasr.Calculate(index - 1)

	var result big.Decimal
	if hasr.isSmooth {
		result = (lastVal.Add(todayVal)).Div(big.NewFromInt(2))
		cacheResult(hasr, index, result)
	} else {
		result = todayVal
		cacheResult(hasr, index, result)
	}

	return result
}

func (i hiekinAshiSmoothedRSIIndicator) cache() resultCache {
	return i.resultCache
}

func (i hiekinAshiSmoothedRSIIndicator) setCache(cache resultCache) {
	i.resultCache = cache
}

func (i hiekinAshiSmoothedRSIIndicator) windowSize() int {
	return i.window
}
