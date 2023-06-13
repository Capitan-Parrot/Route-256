package model

type ClientID uint64
type ClientName string

const (
	ClientUserID   = ClientID(1)
	ClientUserName = ClientName("Gopher")
)

type Client struct {
	ID   ClientID
	Name ClientName
}

type DriverID uint64
type DriverName string
type DriverCar string

const (
	DriverUserID   = DriverID(2)
	DriverUserName = DriverName("Lightning McQueen")
	DriverUserCar  = DriverCar("UAZ Patriot T777TT777")
)

type Driver struct {
	ID   DriverID
	Name DriverName
	Car  DriverCar
}
