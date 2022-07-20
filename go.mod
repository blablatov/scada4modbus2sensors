module github.com/blablatov/scada/main4sensors

go 1.16

require (
	github.com/blablatov/scada/main4sensors/funsensors v0.0.0-20220719183135-c68ec8362197
	github.com/blablatov/scada/main4sensors/gsensors2mgo v0.0.0-00010101000000-000000000000
	github.com/blablatov/scada/main4sensors/sensors2mgo v0.0.0-00010101000000-000000000000
	github.com/blablatov/scada/sensors2mssql v0.0.0-20220719183135-c68ec8362197
)

replace github.com/blablatov/scada/main4sensors/funsensors => ./funsensors

replace github.com/blablatov/scada/main4sensors/gsensors2mgo => ./gsensors2mgo

replace github.com/blablatov/scada/main4sensors/sensors2mgo => ./sensors2mgo
