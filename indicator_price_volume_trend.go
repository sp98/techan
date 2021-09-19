package techan

import "github.com/sdcoffey/big"

type pvtIndicator struct {
	closePriceChangeIndicator Indicator
	volumeIndicator           Indicator
	window                    int
	resultCache               resultCache
}

func NewPriceVolumeTrendIndicator(closePriceIndicator, volumeIndicator Indicator, window int) Indicator {
	return &pvtIndicator{
		closePriceChangeIndicator: NewPercentChangeIndicator(closePriceIndicator),
		volumeIndicator:           volumeIndicator,
		window:                    window,
		resultCache:               make([]*big.Decimal, 1000),
	}
}

func (pvt *pvtIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(pvt, index, func(i int) big.Decimal {
		return big.ZERO
	}); cachedValue != nil {
		return *cachedValue
	}

	priceVolumeChange := pvt.volumeIndicator.Calculate(index).Mul(pvt.closePriceChangeIndicator.Calculate(index))
	previousPVT := pvt.Calculate(index - 1)
	result := priceVolumeChange.Add(previousPVT)

	cacheResult(pvt, index, result)

	return result
}

func (pvt *pvtIndicator) cache() resultCache { return pvt.resultCache }

func (pvt *pvtIndicator) setCache(newCache resultCache) {
	pvt.resultCache = newCache
}

func (pvt *pvtIndicator) windowSize() int { return pvt.window }

func NewPVTAndSignalIndicator(pvtIndicator, signalIndicator Indicator) Indicator {
	return NewDifferenceIndicator(pvtIndicator, signalIndicator)
}
