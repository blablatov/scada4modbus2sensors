module github.com/blablatov/scada4modbus2sensors

go 1.16

require (
	github.com/blablatov/scada4modbus2sensors/funsensors v0.0.0-00010101000000-000000000000
	github.com/blablatov/scada4modbus2sensors/gsensors2mgo v0.0.0-00010101000000-000000000000
	github.com/blablatov/scada4modbus2sensors/sensors2mgo v0.0.0-00010101000000-000000000000
)

replace github.com/blablatov/scada4modbus2sensors/funsensors => ./funsensors

replace github.com/blablatov/scada4modbus2sensors/gsensors2mgo => ./gsensors2mgo

replace github.com/blablatov/scada4modbus2sensors/sensors2mgo => ./sensors2mgo
