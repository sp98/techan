package techan

import (
	"math"

	"github.com/sdcoffey/big"
)

type highestValueIndicator struct {
	indicator Indicator
	window    int
}

func NewHighestValueIndicator(indicator Indicator, window int) Indicator {
	return highestValueIndicator{indicator, window}
}

func (h highestValueIndicator) Calculate(index int) big.Decimal {
	if h.indicator.Calculate(index).NaN() && h.window != 1 {
		return NewHighestValueIndicator(h.indicator, h.window-1).Calculate(index - 1)
	}
	end := math.Max(0, float64(index-h.window+1))
	highest := h.indicator.Calculate(index)
	for i := index - 1; i >= int(end); i-- {
		prevValue := h.indicator.Calculate(i)
		if highest.LT(prevValue) {
			highest = prevValue
		}
	}

	return highest
}
