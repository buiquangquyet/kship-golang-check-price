package ghtk

type GetPriceUnder20Input struct {
	PickAddress  string `json:"pick_address"`
	PickProvince string `json:"pick_province"`
	PickDistrict string `json:"pick_district"`
	PickWard     string `json:"pick_ward"`
	Province     string `json:"province"`
	District     string `json:"district"`
	Ward         string `json:"ward"`
	Address      string `json:"address"`
	Weight       int64  `json:"weight"`
	Value        int64  `json:"value"`
	Transport    string `json:"transport"`
	Tags         []int  `json:"tags"`
}

type GetPriceUnder20Output struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Fee     struct {
		Name         string `json:"name"`
		Fee          int64  `json:"fee"`
		InsuranceFee int64  `json:"insurance_fee"`
		IncludeVat   string `json:"include_vat"`
		CostId       string `json:"cost_id"`
		DeliveryType string `json:"delivery_type"`
		A            string `json:"a"`
		Dt           string `json:"dt"`
		ExtFees      []*Fee `json:"ext_fees"`
		ShipFeeOnly  int64  `json:"ship_fee_only"`
		PromotionKey string `json:"promotion_key"`
		Delivery     bool   `json:"delivery"`
	} `json:"fee"`
}
