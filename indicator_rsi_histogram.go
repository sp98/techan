package techan

import "github.com/sdcoffey/big"

type rsiHistogramIndicator struct {
	rsiIndicator Indicator
	rsiModifer   big.Decimal
}

// NewRSIHistogramIndicator
func NewRSIHistogramIndicator(indicator Indicator, timeframe int) Indicator {
	return rsiHistogramIndicator{
		rsiIndicator: NewRelativeStrengthIndexIndicator(indicator, timeframe),
		rsiModifer:   big.NewFromString("1.5"),
	}
}

func (rsih rsiHistogramIndicator) Calculate(index int) big.Decimal {
	relativeStrengthIndex := rsih.rsiIndicator.Calculate(index)

	return (relativeStrengthIndex.Sub(big.NewFromInt(50))).Mul(rsih.rsiModifer)
}
