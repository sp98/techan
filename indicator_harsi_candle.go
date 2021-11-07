package techan

import (
	"github.com/sdcoffey/big"
)

// Close Price
type harsiClosePriceIndicator struct {
	closeHAIndicator Indicator
	highHAIndicator  Indicator
	lowHAIndicator   Indicator
}

func NewHARSIClosePriceIndicator(series *TimeSeries, timeframe int) Indicator {
	return harsiClosePriceIndicator{
		closeHAIndicator: NewHiekinAshiSmoothedRSIIndicator(NewClosePriceIndicator(series), timeframe, false),
		highHAIndicator:  NewHiekinAshiSmoothedRSIIndicator(NewHighPriceIndicator(series), timeframe, false),
		lowHAIndicator:   NewHiekinAshiSmoothedRSIIndicator(NewLowPriceIndicator(series), timeframe, false),
	}
}

func (h harsiClosePriceIndicator) Calculate(index int) big.Decimal {
	closeRSI := h.closeHAIndicator.Calculate(index)
	openRSI := h.closeHAIndicator.Calculate(index - 1)
	highRSIRAW := h.highHAIndicator.Calculate(index)
	lowRSIRAW := h.lowHAIndicator.Calculate(index)
	highRSI := max(highRSIRAW, lowRSIRAW)
	lowRSI := min(highRSIRAW, lowRSIRAW)
	return (openRSI.Add(highRSI).Add(lowRSI).Add(closeRSI)).Div(big.NewFromInt(4))
}

// Open Price
type harsiOpenPriceIndicator struct {
	closeHAIndicator    Indicator
	closePriceIndicator Indicator
	resultCache         resultCache
	smoothing           big.Decimal
	window              int
}

func NewHARSIOpenPriceIndicator(series *TimeSeries, timeframe int, prevClose Indicator) Indicator {
	return harsiOpenPriceIndicator{
		smoothing:           big.NewFromInt(7),
		closePriceIndicator: prevClose,
		resultCache:         make([]*big.Decimal, 10000),
		window:              timeframe,
		closeHAIndicator:    NewHiekinAshiSmoothedRSIIndicator(NewClosePriceIndicator(series), timeframe, false),
	}
}

func (h harsiOpenPriceIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached2(h, index, func(i int) big.Decimal {
		closeRSI := h.closeHAIndicator.Calculate(index)
		openRSI := h.closeHAIndicator.Calculate(index)
		return (openRSI.Add(closeRSI)).Div(big.NewFromInt(2))
	}); cachedValue != nil {
		return *cachedValue
	}

	result := ((h.Calculate(index - 1).Mul(big.NewFromInt(7))).Add(h.closePriceIndicator.Calculate(index - 1))).Div(big.NewFromInt(8))

	cacheResult(h, index, result)

	return result
}

func (h harsiOpenPriceIndicator) cache() resultCache {
	return h.resultCache
}

func (h harsiOpenPriceIndicator) setCache(cache resultCache) {
	h.resultCache = cache
}

func (h harsiOpenPriceIndicator) windowSize() int {
	return h.window
}

// High Price
type harsiHighPriceIndicator struct {
	openPrice       big.Decimal
	closePrice      big.Decimal
	highHAIndicator Indicator
	lowHAIndicator  Indicator
}

func NewHARSIHighPriceIndicator(open, close big.Decimal, lowHAIndicator, highHAIndicator Indicator) Indicator {
	return harsiHighPriceIndicator{
		openPrice:       open,
		closePrice:      close,
		lowHAIndicator:  lowHAIndicator,
		highHAIndicator: highHAIndicator,
	}
}

func (h harsiHighPriceIndicator) Calculate(index int) big.Decimal {
	highRSIRAW := h.highHAIndicator.Calculate(index)
	lowRSIRAW := h.lowHAIndicator.Calculate(index)
	highRSI := max(highRSIRAW, lowRSIRAW)
	return max(highRSI, max(h.openPrice, h.closePrice))
}

// Low Price
type harsiLowPriceIndicator struct {
	openPrice       big.Decimal
	closePrice      big.Decimal
	highHAIndicator Indicator
	lowHAIndicator  Indicator
}

func NewHARSILowPriceIndicator(open, close big.Decimal, lowHAIndicator, highHAIndicator Indicator) Indicator {
	return harsiLowPriceIndicator{
		openPrice:       open,
		closePrice:      close,
		lowHAIndicator:  lowHAIndicator,
		highHAIndicator: highHAIndicator,
	}
}

func (h harsiLowPriceIndicator) Calculate(index int) big.Decimal {
	highRSIRAW := h.highHAIndicator.Calculate(index)
	lowRSIRAW := h.lowHAIndicator.Calculate(index)
	lowRSI := min(highRSIRAW, lowRSIRAW)
	return min(lowRSI, min(h.openPrice, h.closePrice))
}

func max(a, b big.Decimal) big.Decimal {
	if a.GT(b) {
		return a
	}
	return b
}

func min(a, b big.Decimal) big.Decimal {
	if a.LT(b) {
		return a
	}
	return b
}
