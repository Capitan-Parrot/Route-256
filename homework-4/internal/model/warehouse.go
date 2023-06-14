package model

type WarehouseID uint64
type WarehouseAddress string
type WarehouseWorkHours string

const (
	WarehouseOfficeID        = 2
	WarehouseOfficeAddress   = WarehouseAddress("Saratov")
	WarehouseOfficeWorkHours = WarehouseWorkHours("11-15")
)

/*
Не смог придумать, как полноценно создавать сущность
без клиентов или "загрязнения" order сервера
*/

type Warehouse struct {
	ID        WarehouseID
	Address   WarehouseAddress
	WorkHours WarehouseWorkHours
}
