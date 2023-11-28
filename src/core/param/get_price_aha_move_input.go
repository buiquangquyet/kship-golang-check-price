package param

type GetPriceAhaMoveParam struct {
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
	Requests []*Request
}

type Request struct {
	Id       string
	Num      int
	TierCode string
}
