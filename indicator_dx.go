package techan

import "github.com/sdcoffey/big"

type dxIndicator struct {
	window           int
	plusDIIndicator  Indicator
	minusDIIndicator Indicator
}

func NewDXIndicator(series *TimeSeries, window int) Indicator {
	return dxIndicator{
		window:           window,
		plusDIIndicator:  NewPlusDIIndicator(series, window),
		minusDIIndicator: NewMinusDIIndicator(series, window),
	}
}

func (d dxIndicator) Calculate(index int) big.Decimal {

	plusDIValue := d.plusDIIndicator.Calculate(index)
	minusDIValue := d.minusDIIndicator.Calculate(index)
	if plusDIValue.Add(minusDIValue).EQ(big.NewFromInt(0)) {
		return big.ZERO
	}

	return ((plusDIValue.Sub(minusDIValue).Abs()).Div(plusDIValue.Add(minusDIValue))).Mul(big.NewFromInt(100))
}
