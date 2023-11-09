package constant

const (
	UseKv       = 1
	UseDelivery = 2

	TypeVoucherKv       = 2
	TypeVoucherDelivery = 3

	VoucherExist    = 200
	VoucherNotExist = 999
	VoucherError    = 500
)

var (
	ClientAllowUsePromotion = []string{GRABDeliveryCode, AHAMOVEDeliveryCode, GHNDeliveryCode, GHNFWDeliveryCode}
)
