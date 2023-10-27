package constant

const (
	CodMin = 0
	CodMax = 3000000

	MaxLevel = 9

	ShopVip    = 2
	ShopNormal = 1
	ShopTest   = 0
)

var (
	CodLevelMax = []int{300000, 600000, 1000000, 1200000, 1500000, 1800000, 2200000, 2500000, 3000000}
	CodLevelMin = []int{1, 300001, 600001, 1000001, 1200001, 1500001, 1800001, 2200001, 2500001}

	PriceVip    = []int{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 10000}
	PriceNormal = []int{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 10000}
)
