package constant

const (
	PaymentByFrom = "NGUOIGUI"
	PaymentByTo   = "NGUOINHAN"

	ServiceExtraCodePayment = "PaymentBy"
	ServiceExtraCodeTip     = "TIP"
	ServiceExtraXmg         = "XMG"
	ServiceExtraPickShift   = "PickShift"

	ServiceExtraBulky       = "BULKY"
	ServiceExtraCODST       = "CODST"
	ServiceExtraConn        = "CONN"
	ServiceExtraPartSign    = "partsign"
	ServiceExtraThermalBag  = "THERMALBAG"
	ServiceExtraGbh         = "GBH"
	ServiceExtraShipperNote = "ShipperNote"
	ServiceExtraGnG         = "GNG"
	ServiceExtraRoundTrip   = "ROUND-TRIP"
	ServiceExtraPrepaid     = "PREPAID"
	ServiceExtraBaoPhat     = "BAOPHAT"
	ServiceExtraPtt         = "PTT"
	ServiceExtraDk          = "DK"
	ServiceExtraGHTK1       = "GHTK_1"
	ServiceExtraGHTK7       = "GHTK_7"
	ServiceExtraGHTKXFAST   = "GHTK_XFAST"
	ServiceExtraGHTK2       = "GHTK_2"
	ServiceExtraGHTK3       = "GHTK_3"
	ServiceExtraGHTK4       = "GHTK_4"
	ServiceExtraGHTK5       = "GHTK_5"
)

var (
	MapGHTKExtraService = map[string]int{
		ServiceExtraGHTK1:   1,
		ServiceExtraGHTK7:   7,
		ServiceExtraGHTK2:   13,
		ServiceExtraGHTK3:   17,
		ServiceExtraGHTK4:   18,
		ServiceExtraGHTKTIP: 19,
		ServiceExtraGHTK5:   20,
		ServiceExtraGHTK6:   22,
	}
)
