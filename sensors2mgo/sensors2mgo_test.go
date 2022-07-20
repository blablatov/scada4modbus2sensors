package sensors2mgo

import (
	"fmt"
	"testing"
)

func TestSensData(t *testing.T) {
	var tests = []struct {
		SensorType string
		DataSensor string
	}{
		{"", "\n"},
		{" ", ""},
		{"\t", "one\ttwo\tthree\n"},
		{"Data for test", "#&U*(()))_+_11234"},
		{"Yes, no", "true, false, null"},
	}

	var prevSensorType string
	for _, test := range tests {
		if test.SensorType != prevSensorType {
			fmt.Printf("\n%s\n", test.SensorType)
			prevSensorType = test.SensorType
		}
	}

	var prevDataSensor string
	for _, test := range tests {
		if test.DataSensor != prevDataSensor {
			fmt.Printf("\n%s\n", test.DataSensor)
			prevDataSensor = test.DataSensor
		}
	}
}
