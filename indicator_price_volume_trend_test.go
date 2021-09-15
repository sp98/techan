package techan

import "testing"

func TestPrceVolumeTrendIndicator(t *testing.T) {
	indicator := NewPriceVolumeTrendIndicator(NewClosePriceIndicator(mockedTimeSeries),
		NewVolumeIndicator(mockedTimeSeries), 3)

	expectedValues := []float64{
		0,
		0,
		63.73,
		63.73,
		63.5505,
		63.1925,
		63.9208,
		63.8608,
		62.9735,
		63.3963,
		61.422,
		61.6025,
	}
	indicatorEquals(t, expectedValues, indicator)
}
