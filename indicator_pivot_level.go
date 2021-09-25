package techan

import "github.com/sdcoffey/big"

type PivotLevel int

const (
	RESISTANCE_3 PivotLevel = iota
	RESISTANCE_2
	RESISTANCE_1
	SUPPORT_1
	SUPPORT_2
	SUPPORT_3
)

type pivotLevelIndicator struct {
	indicator  Indicator
	series     *TimeSeries
	pivotLevel PivotLevel
	timeLevel  TimeLevel
}

func NewPivotLevelIndicator(series *TimeSeries, timeLevel TimeLevel, pivotLevel PivotLevel) Indicator {
	return pivotLevelIndicator{
		indicator:  NewPivoPointIndicator(series, timeLevel),
		series:     series,
		pivotLevel: pivotLevel,
		timeLevel:  timeLevel,
	}
}

func (p pivotLevelIndicator) Calculate(index int) big.Decimal {
	pivotPoint := p.indicator.Calculate(index)
	perviousPeriodSeriesIndex := getPreviousPeriodSeries(index, p.series, p.timeLevel)

	switch p.pivotLevel {
	case RESISTANCE_3:
		return calculateR3(pivotPoint, p.series, perviousPeriodSeriesIndex)
	case RESISTANCE_2:
		return calculateR2(pivotPoint, p.series, perviousPeriodSeriesIndex)
	case RESISTANCE_1:
		return calculateR1(pivotPoint, p.series, perviousPeriodSeriesIndex)
	case SUPPORT_3:
		return calculateS3(pivotPoint, p.series, perviousPeriodSeriesIndex)
	case SUPPORT_2:
		return calculateS2(pivotPoint, p.series, perviousPeriodSeriesIndex)
	case SUPPORT_1:
		return calculateS1(pivotPoint, p.series, perviousPeriodSeriesIndex)

	default:
		return big.ZERO
	}
}

func calculateR3(pivotPoint big.Decimal, series *TimeSeries, previousPeriodIndexes []int) big.Decimal {
	lastCandle := series.GetCandle(previousPeriodIndexes[0])
	high := lastCandle.MaxPrice
	low := lastCandle.MinPrice

	for _, index := range previousPeriodIndexes {
		nextHigh := series.GetCandle(previousPeriodIndexes[index]).MaxPrice
		if nextHigh.GT(high) {
			high = nextHigh
		}
		nextLow := series.GetCandle(previousPeriodIndexes[index]).MinPrice
		if nextLow.LT(low) {
			low = nextLow
		}
	}

	return high.Add(big.NewFromInt(2).Mul((pivotPoint.Sub(low))))
}

func calculateR2(pivotPoint big.Decimal, series *TimeSeries, previousPeriodIndexes []int) big.Decimal {
	lastCandle := series.GetCandle(previousPeriodIndexes[0])
	high := lastCandle.MaxPrice
	low := lastCandle.MinPrice

	for _, index := range previousPeriodIndexes {
		nextHigh := series.GetCandle(previousPeriodIndexes[index]).MaxPrice
		if nextHigh.GT(high) {
			high = nextHigh
		}
		nextLow := series.GetCandle(previousPeriodIndexes[index]).MinPrice
		if nextLow.LT(low) {
			low = nextLow
		}
	}

	return pivotPoint.Add((high.Sub(low)))
}

func calculateR1(pivotPoint big.Decimal, series *TimeSeries, previousPeriodIndexes []int) big.Decimal {
	lastCandle := series.GetCandle(previousPeriodIndexes[0])
	low := lastCandle.MinPrice

	for _, index := range previousPeriodIndexes {
		nextLow := series.GetCandle(previousPeriodIndexes[index]).MinPrice
		if nextLow.LT(low) {
			low = nextLow
		}
	}

	return (big.NewFromInt(2).Mul(pivotPoint).Sub(low))
}

func calculateS1(pivotPoint big.Decimal, series *TimeSeries, previousPeriodIndexes []int) big.Decimal {
	lastCandle := series.GetCandle(previousPeriodIndexes[0])
	high := lastCandle.MaxPrice

	for _, index := range previousPeriodIndexes {
		nextHigh := series.GetCandle(previousPeriodIndexes[index]).MaxPrice
		if nextHigh.GT(high) {
			high = nextHigh
		}
	}

	return (big.NewFromInt(2).Mul(pivotPoint).Sub(high))
}

func calculateS2(pivotPoint big.Decimal, series *TimeSeries, previousPeriodIndexes []int) big.Decimal {
	lastCandle := series.GetCandle(previousPeriodIndexes[0])
	high := lastCandle.MaxPrice
	low := lastCandle.MinPrice

	for _, index := range previousPeriodIndexes {
		nextHigh := series.GetCandle(previousPeriodIndexes[index]).MaxPrice
		if nextHigh.GT(high) {
			high = nextHigh
		}
		nextLow := series.GetCandle(previousPeriodIndexes[index]).MinPrice
		if nextLow.LT(low) {
			low = nextLow
		}
	}

	return pivotPoint.Sub((high.Sub(low)))
}

func calculateS3(pivotPoint big.Decimal, series *TimeSeries, previousPeriodIndexes []int) big.Decimal {
	lastCandle := series.GetCandle(previousPeriodIndexes[0])
	high := lastCandle.MaxPrice
	low := lastCandle.MinPrice

	for _, index := range previousPeriodIndexes {
		nextHigh := series.GetCandle(previousPeriodIndexes[index]).MaxPrice
		if nextHigh.GT(high) {
			high = nextHigh
		}
		nextLow := series.GetCandle(previousPeriodIndexes[index]).MinPrice
		if nextLow.LT(low) {
			low = nextLow
		}
	}

	return low.Sub(big.NewFromInt(2).Mul((high.Sub(pivotPoint))))
}
