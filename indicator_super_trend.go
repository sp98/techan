package techan

import "github.com/sdcoffey/big"

type superTrendIndicator struct {
	upperBandIndicator  Indicator
	lowerBandIndicator  Indicator
	closePriceIndicator Indicator
	resultCache         resultCache
	window              int
}

func NewSuperTrendIndicator(series *TimeSeries, window, atrLength, multiplier int) Indicator {
	atrIndicator := NewAverageTrueRangeIndicator(series, atrLength)
	return superTrendIndicator{
		upperBandIndicator:  NewUpperBandIndicator(series, atrIndicator, window, multiplier),
		lowerBandIndicator:  NewLowerBandIndicator(series, atrIndicator, window, multiplier),
		closePriceIndicator: NewClosePriceIndicator(series),
		window:              window,
	}
}

func (st superTrendIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(st, index, func(i int) big.Decimal {
		return big.ZERO
	}); cachedValue != nil {
		return *cachedValue
	}

	stPrevious := st.Calculate(index - 1)
	ubCurrent := st.upperBandIndicator.Calculate(index)
	ubPrevious := st.upperBandIndicator.Calculate(index - 1)
	lbCurrent := st.lowerBandIndicator.Calculate(index)
	lbPrevious := st.lowerBandIndicator.Calculate(index - 1)
	var result big.Decimal

	if stPrevious.EQ(ubPrevious) && st.closePriceIndicator.Calculate(index).LTE(ubCurrent) {
		result = ubCurrent
	} else if stPrevious.EQ(ubPrevious) && st.closePriceIndicator.Calculate(index).GTE(ubCurrent) {
		result = lbCurrent
	} else if stPrevious.EQ(lbPrevious) && st.closePriceIndicator.Calculate(index).GTE(lbCurrent) {
		result = lbCurrent
	} else if stPrevious.EQ(lbPrevious) && st.closePriceIndicator.Calculate(index).LTE(lbCurrent) {
		result = ubCurrent
	} else {
		result = big.ZERO
	}

	cacheResult(st, index, result)

	return result
}

func (st superTrendIndicator) cache() resultCache { return st.resultCache }

func (st superTrendIndicator) setCache(newCache resultCache) {
	st.resultCache = newCache
}

func (st superTrendIndicator) windowSize() int { return st.window }

/* Upper and Lower Band */
type UpperBandIndicator struct {
	closePriceIndicator    Indicator
	upperBandBasicIndictor Indicator
	window                 int
	resultCache            resultCache
}

func NewUpperBandIndicator(series *TimeSeries, artIndicator Indicator, window, multiplier int) Indicator {
	return &UpperBandIndicator{
		closePriceIndicator:    NewClosePriceIndicator(series),
		upperBandBasicIndictor: NewUpperBandBasicIndicator(series, artIndicator, multiplier),
		window:                 window,
	}
}

func (ub UpperBandIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(ub, index, func(i int) big.Decimal {
		return big.ZERO
	}); cachedValue != nil {
		return *cachedValue
	}
	ubPrevious := ub.Calculate(index - 1)
	ubbCurrent := ub.upperBandBasicIndictor.Calculate(index)
	ubbPrevious := ub.upperBandBasicIndictor.Calculate(index - 1)
	if ubbCurrent.LT(ubbPrevious) || ub.closePriceIndicator.Calculate(index-1).GT(ubPrevious) {
		return ubbCurrent
	}
	return ubbPrevious
}

func (ub UpperBandIndicator) cache() resultCache { return ub.resultCache }

func (ub UpperBandIndicator) setCache(newCache resultCache) {
	ub.resultCache = newCache
}

func (ub UpperBandIndicator) windowSize() int { return ub.window }

type LowerBandIndicator struct {
	closePriceIndicator    Indicator
	lowerBandBasicIndictor Indicator
	window                 int
	resultCache            resultCache
}

func NewLowerBandIndicator(series *TimeSeries, artIndicator Indicator, window, multiplier int) Indicator {
	return &LowerBandIndicator{
		closePriceIndicator:    NewClosePriceIndicator(series),
		lowerBandBasicIndictor: NewLowerBandBasicIndicator(series, artIndicator, multiplier),
		window:                 window,
	}
}

func (lb LowerBandIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(lb, index, func(i int) big.Decimal {
		return big.ZERO
	}); cachedValue != nil {
		return *cachedValue
	}
	lbPrevious := lb.Calculate(index - 1)
	lbbCurrent := lb.lowerBandBasicIndictor.Calculate(index)
	lbbPrevious := lb.lowerBandBasicIndictor.Calculate(index - 1)
	var result big.Decimal
	if lbbCurrent.GT(lbbPrevious) || lb.closePriceIndicator.Calculate(index-1).LT(lbPrevious) {
		result = lbbCurrent
	} else {
		result = lbPrevious
	}

	cacheResult(lb, index, result)

	return result
}

func (lb LowerBandIndicator) cache() resultCache { return lb.resultCache }

func (lb LowerBandIndicator) setCache(newCache resultCache) {
	lb.resultCache = newCache
}

func (lb LowerBandIndicator) windowSize() int { return lb.window }

/* Basic Band Calculation */

type UpperBandBasicIndicator struct {
	highPriceIndicator Indicator
	lowPriceIndicator  Indicator
	atrIndicator       Indicator
	multiplier         int
	// window             int
	// resultCache        resultCache
}

func NewUpperBandBasicIndicator(series *TimeSeries, artIndicator Indicator, multiplier int) Indicator {
	return &UpperBandBasicIndicator{
		highPriceIndicator: NewHighPriceIndicator(series),
		lowPriceIndicator:  NewLowPriceIndicator(series),
		atrIndicator:       artIndicator,
		multiplier:         multiplier,
	}
}

func (ubb UpperBandBasicIndicator) Calculate(index int) big.Decimal {
	// if cachedValue := returnIfCached(ubb, index, func(i int) big.Decimal {
	// 	return big.ZERO
	// }); cachedValue != nil {
	// 	return *cachedValue
	// }
	if index == 0 {
		return big.ZERO
	}

	return (ubb.highPriceIndicator.Calculate(index).Add(ubb.lowPriceIndicator.Calculate(index).Div(big.NewFromInt(2)))).Sub((ubb.atrIndicator.Calculate(index)).Mul(big.NewFromInt(ubb.multiplier)))
}

// func (st UpperBandBasicIndicator) cache() resultCache { return st.resultCache }

// func (st UpperBandBasicIndicator) setCache(newCache resultCache) {
// 	st.resultCache = newCache
// }

// func (st UpperBandBasicIndicator) windowSize() int { return st.window }

type LowerBandBasicIndicator struct {
	highPriceIndicator Indicator
	lowPriceIndicator  Indicator
	atrIndicator       Indicator
	multiplier         int
}

func NewLowerBandBasicIndicator(series *TimeSeries, artIndicator Indicator, multiplier int) Indicator {
	return &LowerBandBasicIndicator{
		highPriceIndicator: NewHighPriceIndicator(series),
		lowPriceIndicator:  NewLowPriceIndicator(series),
		atrIndicator:       artIndicator,
		multiplier:         multiplier,
	}
}

func (lbb LowerBandBasicIndicator) Calculate(index int) big.Decimal {
	if index == 0 {
		return big.ZERO
	}

	return (lbb.highPriceIndicator.Calculate(index).Sub(lbb.lowPriceIndicator.Calculate(index).Div(big.NewFromInt(2)))).Sub((lbb.atrIndicator.Calculate(index)).Mul(big.NewFromInt(lbb.multiplier)))
}
