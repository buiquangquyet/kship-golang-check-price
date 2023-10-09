package dto

type GetPriceInputDto struct {
	PickProvince     string
	PickDistrict     string
	PickWard         string
	ReceiverProvince string
	ReceiverDistrict string
	ReceiverWard     string
	Address          string
	Products         []*Product //BBS
	Weight           int64      //no BBS
	Value            int64
	Transport        string
	OrderService     string
	Tags             []int
}

type Product struct {
	Name     string
	Quantity int
	Weight   int64
	Width    int
	Length   int
	Height   int
}