package main

import (
	"fmt"
	"gsensors2mgo"
	"sensors2mgo"
	"sensors2mssql"
	"sync"
	"testing"
)

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
	for i := 0; i < 5; i++ {
		md := sensors2mgo.SensMongo{
			SensorType: "dallas_1",
			DataSensor: 18.433,
		}
		ct := make(chan string)   // Channel to data send  type of sensor. Канал передачи типа датчика.
		ctd := make(chan float64) // Channel to data send temperature. Канал данных температуры.
		//done := make(chan bool)   // Channel of synchronization. Канал синхронизации.
		var wg sync.WaitGroup
		go gsensors2mgo.SensData(md.SensorType, DsnMongo, md.DataSensor, ct, ctd, wg)
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
