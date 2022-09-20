module github.com/blablatov/scada4modbus2sensors

go 1.16

replace github.com/blablatov/scada4modbus2sensors/funsensors => ./funsensors

replace github.com/blablatov/scada4modbus2sensors/gsensors2mgo => ./gsensors2mgo

replace github.com/blablatov/scada4modbus2sensors/sensors2mgo => ./sensors2mgo

require (
	github.com/blablatov/scada4modbus2sensors/funsensors v0.0.0-20220920103001-40558a815aa8
	github.com/blablatov/scada4modbus2sensors/gsensors2mgo v0.0.0-20220920103001-40558a815aa8
	github.com/blablatov/scada4modbus2sensors/sensors2mgo v0.0.0-20220920103001-40558a815aa8
	github.com/blablatov/scada4sensors2mssql v0.0.0-20220920115510-38baea04e421
)
