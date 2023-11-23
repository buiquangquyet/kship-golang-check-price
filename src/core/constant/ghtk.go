package constant

var (
	MapGHTKTagByServiceExtra = map[string]int{
		ServiceExtraGHTK1:   1,
		ServiceExtraGHTK7:   7,
		ServiceExtraGHTK2:   13,
		ServiceExtraGHTK3:   17,
		ServiceExtraGHTK4:   18,
		ServiceExtraCodeTip: 19,
		ServiceExtraGHTK5:   20,
		ServiceExtraGHTK6:   22,
	}
)

var (
	MapGHTKTagByShipperNote = map[string]int{
		"CHOXEMHANGKHONGTHU": 10,
		"CHOTHUHANG":         11,
	}
)

var (
	MinFailDeliveryTaking int64 = 10000
	MaxFailDeliveryTaking int64 = 20000000
)
