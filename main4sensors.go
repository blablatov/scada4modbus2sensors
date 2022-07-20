// Main demo module do formating and transmit data to MongoDB or MSSQL module for write to DB.
// Основной демо модуль выполняет форматирование и передает данные в модуль MongoDB или MSSQL для записи в БД.
// Base of idea to https://github.com/samkalnins/ds18b20-prometheus-exporter. Thanks him!
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/blablatov/scada/sensors2mssql"

	"github.com/blablatov/scada4modbus2sensors/funsensors"
	"github.com/blablatov/scada4modbus2sensors/gsensors2mgo"
	"github.com/blablatov/scada4modbus2sensors/sensors2mgo"
)

type dBaseMssql struct {
	dBaseName string //`json:"dbserver"`
	//dBaseType   string `json:"dbserver"`
}

const (
	DsnMongo = "mongodb://localhost:27017/testdb"
)

var bus_dir_t = flag.String("w1_bus_dir", "./w1_slave", "directory of the 1-wire bus")

// For another sensors, pressure and etc.
//var bus_dir_p = flag.String("w1_bus_dir", "/sys/bus/w1/devices", "directory of the 1-wire bus")
//var bus_dir_f = flag.String("w1_bus_dir", "/sys/bus/w1/devices", "directory of the 1-wire bus")

//var port = flag.Int("port", 8008, "port to run http server on")
var port = flag.Int("port", 8443, "port to run http server on")

type prometheusLabels map[string][]string

// String is the method to format the flag's value, part of the flag.Value interface.
func (p *prometheusLabels) String() string {
	return fmt.Sprint(*p)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (p *prometheusLabels) Set(value string) error {
	*p = make(map[string][]string)

	for _, ls := range strings.Split(value, ",") {
		s := strings.Split(ls, "=")
		if len(s) != 3 {
			errors.New("Bad flag value -- should be temp_id=label=value")
		}
		_, initialized := (*p)[s[0]]
		if !initialized {
			(*p)[s[0]] = make([]string, 0)
		}
		(*p)[s[0]] = append((*p)[s[0]], fmt.Sprintf("%s=\"%s\"", s[1], s[2]))
	}
	return nil
}

var prometheusLabelsFlag prometheusLabels

func init() {
	//При наличии нескольких датчиков, опрос выполняется по номеру датчика связанному с параметром - локация.
	//$ temp_exporter --port 8000 \ --prometheus_labels "28-0416a4a474ff=location=lounge,"28-0417713760ff"=location=garden"
	flag.Var(&prometheusLabelsFlag, "prometheus_labels", "comma-separated list of labels to apply to sensors by ID e.g. 28-0417713760f=label=value,")
}

func main() {
	flag.Parse()

	// Main varz handler -- read and parse the temperatures on each request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is listening on 8443. Go to https://127.0.0.1:8443")
		// Here reading data from local file, demo it, without sensors.
		// Чтение данных из тестового файла, для отладки, без подключенных датчиков.
		readings_t, err := funsensors.ReadTemperatureFile(*bus_dir_t)
		if err != nil {
			log.Printf("Error reading temperatures [%s]", err)
		}
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
		fmt.Fprintf(w, "Data of sensor: %f\n", readings_t)
		fmt.Printf("Data of sensor for write: %f\n", readings_t)

		////////////////////////////////////////////////////////////////////////
		// Sending data to MongoDB via method of interface.
		// Отправка данных в MongoDB через вызов метода интерфейса.
		start := time.Now()
		// Formating data of structure. Заполнение структуры.
		md := sensors2mgo.SensMongo{
			SensorType: "dallas_1",
			DataSensor: readings_t,
		}
		// Вызов метода интерфейса.
		// Calling an interface method.
		var d sensors2mgo.Monger = md
		db := d.SensData(DsnMongo)
		fmt.Println("Result of request to MongoDB via interface method: ", db)
		secs := time.Since(start).Seconds()
		fmt.Printf("%.2fs Request execution time via method of interface\n", secs)

		////////////////////////////////////////////////////////////////////////
		// Sending data to MongoDB via goroutine.
		// Отправка данных в MongoDB через горутину.
		start2 := time.Now()
		// Formating data of structure. Заполнение структуры.
		mg := gsensors2mgo.SensMongo{
			SensorType: "dallas_2",
			DataSensor: readings_t,
		}
		ct := make(chan string)   // Channel to data send type of sensor. Канал передачи типа датчика.
		ctd := make(chan float64) // Channel to data send temperature. Канал данных температуры.
		//done := make(chan bool)   // Channel of synchronization. Канал синхронизации.
		var wg sync.WaitGroup // Synchronization of goroutines. Синхронизация горутин.
		wg.Add(1)             // Counter of goroutines. Значение счетчика горутин
		go gsensors2mgo.SensData(mg.SensorType, DsnMongo, mg.DataSensor, ct, ctd, wg)
		// Getting data from goroutine. Получение данных из канала горутины.
		log.Println("\nSensor of system: ", <-ct, "\nData of sensor: ", <-ctd)
		// Wait of counter. Ожидание счетчика
		go func() {
			wg.Wait() // Waiting of counter. Ожидание счетчика.
			close(ct)
			close(ctd)
		}()
		secs2 := time.Since(start2).Seconds()
		fmt.Printf("%.2fs Request execution time to MongoDB via goroutine\n", secs2)

		////////////////////////////////////////////////////////////////////////
		// Sending data to MSSQL via method of interface.
		// Отправка данных в MSSQL через вызов метода интерфейса.
		start3 := time.Now()
		// Formating data of structure. Заполнение структуры.
		ds := dBaseMssql{
			dBaseName: "mssqlserver",
			//dBaseType: anydata,
		}
		// Formating data structure of operators. Заполнение структуры операторов.
		rq := sensors2mssql.ReqOperators{
			InsertReqSql: "INSERT " + ds.dBaseName + " VALUES (SensorType, " + md.SensorType + ")",
			SelectReqSql: "SELECT Name FROM TableDB WHERE BusUnit = 65;",
			UpdateReqSql: "UPDATE TableDB SET Id=Id, Name=Name, BusUnit=BusUnit, ItemNumber=ItemNumber, ItemName=ItemName SELECT Id, Name, BusUnit, ItemNumber, ItemName FROM AnyTable WHERE BusUnit = 65",
			Create:       "CREATE TABLE MyTable(Id int, Name nvarchar(250), BusUnit int, ItemNumber nvarchar(50), ItemName nvarchar(100));",
		}
		// Вызов метода интерфейса.
		// Calling an interface method.
		var dm sensors2mssql.ReqOperators = rq
		mr := dm.SensMssql()
		fmt.Println("Result of request to MSSQL via interface method: ", mr)
		secs3 := time.Since(start3).Seconds()
		fmt.Printf("%.2fs Request execution time\n", secs3)

		// Чтение данных из файла драйвера датчика для unix
		// Here reading data from local file with connected sensor.
		/*readings, err := dstemp.FindAndReadTemperatures(*bus_dir)
		if err != nil {
			log.Printf("Error reading temperatures [%s]", err)
		}
		for _, tr := range readings {
			labels := strings.Join(append(prometheusLabelsFlag[tr.Id], fmt.Sprintf("sensor=\"%s\"", tr.Id)), ",")

			// Output varz as both C & F for maximum user happiness
			fmt.Fprintf(w, "temperature_c{%s} %f\n", labels, tr.Temp_c)
			fmt.Fprintf(w, "temperature_f{%s} %f\n", labels, dstemp.CentigradeToF(tr.Temp_c))
		}*/
	})
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), "gserver.crt", "gserver.key", nil))
}
