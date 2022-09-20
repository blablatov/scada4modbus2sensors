package gsensors2mgo

import (
	"fmt"
	"log"
	"testing"

	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

const (
	DsnMongo = "mongodb://localhost:27017/testdb"
)

func TestSensData(t *testing.T) {
	var tests = []struct {
		SensorType string
		DataSensor float64
	}{
		{"dallas_1", 55.6},
		{"dallas_2", -18.2},
		{"dallas_54", -45.1},
		{"Data for test", 0.1},
		{"Yes, no", 0.02},
	}

	var prevSensorType string
	for _, test := range tests {
		if test.SensorType != prevSensorType {
			fmt.Printf("\n%s\n", test.SensorType)
			prevSensorType = test.SensorType
		}
	}

	var prevDataSensor float64
	for _, test := range tests {
		if test.DataSensor != prevDataSensor {
			fmt.Printf("\n%f\n", test.DataSensor)
			prevDataSensor = test.DataSensor
		}
	}

	strDataSensor := fmt.Sprint(prevDataSensor)

	session, err := mgo.Dial(DsnMongo)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// is check name in dBase
	c := session.DB("scadadb").C("sensors")
	chk := SensMongo{}
	err = c.Find(bson.M{"sensortype": prevSensorType, "datasensor": strDataSensor}).One(&chk)
	if err == nil {
		log.Println("\nName already is to DB, method of interface", err)
	}
	if err != nil {
		log.Print("\nErr data for write to MongoDB, method of interface: ", err)
		err = c.Insert(&SensMongo{prevSensorType, prevDataSensor})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sensor was written via method of interface:", prevSensorType,
			"\nData of sensor was written via method of interface:", strDataSensor)
	}
}
