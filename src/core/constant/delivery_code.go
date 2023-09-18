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
)

var (
	ReceiverWardIdClientCode = []string{AHAMOVEDeliveryCode, VNPOSTDeliveryCode, GRABDeliveryCode,
		EMSDeliveryCode, VTPFWDeliveryCode, VIETTELPostDeliveryCode, GHTKDeliveryCode}
)
