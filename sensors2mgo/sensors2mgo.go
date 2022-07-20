package sensors2mgo

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Monger interface {
	SensData(DsnMongo string) bool
}

type SensMongo struct {
	SensorType string
	DataSensor float64
	//Sensor_p string
	//Pressure string
	//Sensor_f string
	//Flow     string
	//Sensor_g string
	//Gas      string
}

func (md SensMongo) SensData(DsnMongo string) bool {
	//session, err := mgo.Dial("mongodb://localhost:27017/testdb")
	session, err := mgo.Dial(DsnMongo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// is check name in dBase
	c := session.DB("scadadb").C("sensors")
	chk := SensMongo{}
	err = c.Find(bson.M{"sensortype": md.SensorType, "datasensor": md.DataSensor}).One(&chk)
	if err == nil {
		log.Print("\nName already is to MongoDB via interface method", err)
		return false
	}
	if err != nil {
		log.Print("\nData in MongoDB via interface method:", err)
		err = c.Insert(&SensMongo{md.SensorType, md.DataSensor})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sensor was written to MongoDB via interface method: ", md.SensorType,
			"\nData was written to MongoDB via interface method: ", md.DataSensor)
	}
	return true
}
