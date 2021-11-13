package techan

import (
	"math"

	"github.com/sdcoffey/big"
)

type lowestValueIndicator struct {
	indicator Indicator
	window    int
}

func NewLowestValueIndicator(indicator Indicator, window int) Indicator {
	return lowestValueIndicator{indicator, window}
}

func (l lowestValueIndicator) Calculate(index int) big.Decimal {
	if l.indicator.Calculate(index).NaN() && l.window != 1 {
		return NewLowestValueIndicator(l.indicator, l.window-1).Calculate(index - 1)
	}
	end := math.Max(0, float64(index-l.window+1))
	lowest := l.indicator.Calculate(index)
	for i := index - 1; i >= int(end); i-- {
		prevValue := l.indicator.Calculate(i)
		if lowest.GT(prevValue) {
			lowest = prevValue
		}
	}

	return lowest
}
