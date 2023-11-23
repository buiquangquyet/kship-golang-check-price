package dto

type GetPriceInputAhaMoveDto struct {
	Path          [2]*Path
	PaymentMethod string
	PromoCode     string
	OrderTime     int64
	Services      []*ServiceAhaMove
}
type Path struct {
	Address string
	Cod     int64
}

type ServiceAhaMove struct {
	Id       string
	Requests []string
}
