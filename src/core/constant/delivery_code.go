package constant

const (
	VTPFWDeliveryCode  = "VTPFW"
	VNPDeliveryCode    = "VNP"
	GHTKDeliveryCode   = "GHTK"
	JTFWDeliveryCode   = "JTFW"
	GHNFWDeliveryCode  = "GHNFW"
	BESTFWDeliveryCode = "BESTFW"

	AHAMOVEDeliveryCode     = "AHAMOVE"
	VNPOSTDeliveryCode      = "VNP"
	GRABDeliveryCode        = "GRAB"
	EMSDeliveryCode         = "EMS"
	VIETTELPostDeliveryCode = "VTP"

	GHNDeliveryCode = "GHN"
	SpxDeliveryCode = "SPX"
)

var (
	SenderWardIdDeliveryCode   = []string{GHTKDeliveryCode}
	ReceiverWardIdDeliveryCode = []string{AHAMOVEDeliveryCode, VNPOSTDeliveryCode, GRABDeliveryCode,
		EMSDeliveryCode, VTPFWDeliveryCode, VIETTELPostDeliveryCode, GHTKDeliveryCode}

	PickShiftExtraClientCode = []string{GHNDeliveryCode, GHNFWDeliveryCode, SpxDeliveryCode}
)
