package model

type PickUpPointID uint64
type PickUpPointAddress string
type PickUpPointWorkHours string

const (
	PickUpPointOfficeID        = 2
	PickUpPointOfficeAddress   = PickUpPointAddress("Saratov")
	PickUpPointOfficeWorkHours = PickUpPointWorkHours("11-15")
)

/*
Не смог придумать, как полноценно создавать сущность
без клиентов или "загрязнения" order сервера
*/

type PickUpPoint struct {
	ID        PickUpPointID
	Address   PickUpPointAddress
	WorkHours PickUpPointWorkHours
}
