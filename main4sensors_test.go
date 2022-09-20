package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/blablatov/scada4modbus2sensors/gsensors2mgo"
	"github.com/blablatov/scada4modbus2sensors/sensors2mgo"
	"github.com/blablatov/scada4sensors2mssql"
)

func Test(t *testing.T) {
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
}

func BenchmarkInterface(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 5; i++ {
		md := sensors2mgo.SensMongo{
			SensorType: "dallas_1",
			DataSensor: 18.433,
		}
		var d sensors2mgo.Monger = md
		db := d.SensData(DsnMongo)
		fmt.Println("Result of request to MongoDB via interface method", db)
	}
}

func BenchmarkGoroutine(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 1; i++ {
		SensorType := "dallas_1"
		DataSensor := 18.433
		ct := make(chan string)   // Channel to data send  type of sensor. Канал передачи типа датчика.
		ctd := make(chan float64) // Channel to data send temperature. Канал данных температуры.
		var wg sync.WaitGroup
		wg.Add(1)
		go gsensors2mgo.SensData(SensorType, DsnMongo, DataSensor, ct, ctd, wg)
		// Getting data from goroutine. Получение данных из канала горутины.
		fmt.Println("\nSensor of system: ", <-ct, "\nData of sensor: ", <-ctd)
		go func() {
			wg.Wait()
			close(ct)
			close(ctd)
		}()
	}
}

func BenchmarkInterfaceMssql(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 5; i++ {
		rq := sensors2mssql.ReqOperators{
			InsertReqSql: "INSERT Name VALUES (SensorType, DataSensor)",
			SelectReqSql: "SELECT Name FROM TableDB WHERE BusUnit = 65;",
			UpdateReqSql: "UPDATE TableDB SET Id=Id, Name=Name, BusUnit=BusUnit, ItemNumber=ItemNumber, ItemName=ItemName SELECT Id, Name, BusUnit, ItemNumber, ItemName FROM AnyTable WHERE BusUnit = 65",
			Create:       "CREATE TABLE MyTable(Id int, Name nvarchar(250), BusUnit int, ItemNumber nvarchar(50), ItemName nvarchar(100));",
		}
		var dm sensors2mssql.ReqOperators = rq
		mr := dm.SensMssql()
		fmt.Println("Result of request to MSSQL via interface method: ", mr)
	}
}
