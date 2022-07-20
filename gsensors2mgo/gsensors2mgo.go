package gsensors2mgo

import (
	"fmt"
	"log"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

func SensData(SensorType, DsnMongo string, DataSensor float64, ct chan string, ctd chan float64, wg sync.WaitGroup) {
	defer wg.Done()
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
	err = c.Find(bson.M{"sensortype": SensorType, "datasensor": DataSensor}).One(&chk)
	if err == nil {
		log.Println("\nName already is to DB via goroutine", err)
		return
	}
	if err != nil {
		log.Print("\nData in MongoDB via goroutine: ", err)
		err = c.Insert(&SensMongo{SensorType, DataSensor})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sensor was written via goroutine:", SensorType, "\nData was written via goroutine:", DataSensor)
		ct <- SensorType
		ctd <- DataSensor
	}
}
