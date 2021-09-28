package techan

import (
	"github.com/sdcoffey/big"
)

type TimeLevel int

const (
	BARBASED TimeLevel = iota
	DAY
	WEEK
	MONTH
)

type pivotPointIndicator struct {
	series    *TimeSeries
	timeLevel TimeLevel
}

func NewPivoPointIndicator(series *TimeSeries, timeLevel TimeLevel) Indicator {
	return pivotPointIndicator{series, timeLevel}
}

func (p pivotPointIndicator) Calculate(index int) big.Decimal {
	return CalculatePivotPoint(p.series, getPreviousPeriodSeries(index, p.series, p.timeLevel))
}

func CalculatePivotPoint(series *TimeSeries, previousPeriodIndexes []int) big.Decimal {
	if len(previousPeriodIndexes) == 0 {
		panic("previous period candles not available")
	}

	lastCandle := series.GetCandle(previousPeriodIndexes[0])
	close := lastCandle.ClosePrice
	high := lastCandle.MaxPrice
	low := lastCandle.MinPrice

	for _, index := range previousPeriodIndexes {
		nextHigh := series.GetCandle(index).MaxPrice
		if nextHigh.GT(high) {
			high = nextHigh
		}
		nextLow := series.GetCandle(index).MinPrice
		if nextLow.LT(low) {
			low = nextLow
		}
	}

	return high.Add(low).Add(close).Div(big.NewFromInt(3))
}

func getPreviousPeriodSeries(index int, series *TimeSeries, timeLevel TimeLevel) []int {

	previousSeriesIndexes := []int{}

	if timeLevel == BARBASED {
		previousSeriesIndexes = append(previousSeriesIndexes, Max(0, index-1))
	}

	if index == 0 {
		return previousSeriesIndexes
	}

	currentCandle := series.GetCandle(index)
	currentCandlePeriod := getPeriod(currentCandle, timeLevel)

	for index-1 > 0 && getPeriod(series.GetCandle(index-1), timeLevel) == currentCandlePeriod {
		index--
	}

	previousPeriod := getPreviousPeriod(index, timeLevel, series)

	for index-1 > 0 && getPeriod(series.GetCandle(index-1), timeLevel) == previousPeriod {
		index--
		previousSeriesIndexes = append(previousSeriesIndexes, index)
	}

	return previousSeriesIndexes
}

func getPeriod(candle *Candle, timeLevel TimeLevel) int {
	switch timeLevel {
	case DAY:
		return candle.Period.Start.YearDay()

	case WEEK:
		_, w := candle.Period.Start.ISOWeek()
		return w

	case MONTH:
		return int(candle.Period.Start.Month())

	default:
		return candle.Period.Start.Year()
	}
}

func getPreviousPeriod(index int, timeLevel TimeLevel, series *TimeSeries) int {
	currentCandle := series.GetCandle(index)
	switch timeLevel {
	case DAY:
		previousCandle := series.GetCandle(index - 1)
		prevCalendarDay := currentCandle.Period.Start.AddDate(0, 0, -1).YearDay()
		// skip weekend and holidays
		for previousCandle.Period.Start.YearDay() != prevCalendarDay && prevCalendarDay >= 0 {
			prevCalendarDay--
		}
		return prevCalendarDay

	case WEEK:
		_, w := currentCandle.Period.Start.ISOWeek()
		return w - 1

	case MONTH:
		return int(currentCandle.Period.Start.AddDate(0, -1, 0).Month())

	default:
		return currentCandle.Period.Start.AddDate(-1, 0, 0).Year()
	}
}
