package param

type GetPriceGHTKParam struct {
	PickProvince     string
	PickDistrict     string
	PickWard         string
	ReceiverProvince string
	ReceiverDistrict string
	ReceiverWard     string
	Address          string
	Products         []*ProductGHTK //BBS
	Weight           int64          //no BBS
	Value            int64
	Tags             []int
	NotDeliveredFee  int64
}

type ProductGHTK struct {
	Name     string
	Quantity int
	Weight   int64
	Width    int
	Length   int
	Height   int
}
