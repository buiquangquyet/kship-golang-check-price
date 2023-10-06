package ghtk

type GetPriceOver20Input struct {
	PickAddress      string     `json:"pick_address"`
	PickProvince     string     `json:"pick_province"`
	PickDistrict     string     `json:"pick_district"`
	PickWard         string     `json:"pick_ward"`
	Value            int64      `json:"value"`
	PickStreet       string     `json:"pick_street"`
	CustomerProvince string     `json:"customer_province"`
	CustomerDistrict string     `json:"customer_district"`
	CustomerWard     string     `json:"customer_ward"`
	CustomerAddress  string     `json:"customer_address"`
	CustomerHamlet   string     `json:"customer_hamlet"`
	Products         []*Product `json:"products"`
	Transport        string     `json:"transport"`
	Tags             []int      `json:"tags"`
}

type Product struct {
	Name     string `json:"name"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Length   int    `json:"length"`
	Quantity int    `json:"quantity"`
	Weight   int    `json:"weight"`
}

type GetPriceOver20Output struct {
	Success bool `json:"success"`
	Data    struct {
		CostId     int64   `json:"cost_id"`
		RealWeight int64   `json:"real_weight"`
		Distance   float64 `json:"distance"`
		Value      int64   `json:"value"`
		Transport  string  `json:"transport"`
		Flag       struct {
			BaseCost map[string]int64 `json:"base_cost"`
			Step     int              `json:"step"`
			Increase int64            `json:"increase"`
		} `json:"flag"`
		Region      string `json:"region"`
		OldValue    int64  `json:"old_value"`
		OnlyShipFee int64  `json:"only_ship_fee"`
		Insurance   int64  `json:"insurance"`
		ExtFees     []*Fee `json:"ext_fees"`
		FragileFee  int64  `json:"fragile_fee"`
		FoodFee     int64  `json:"food_fee"`
		TotalValue  int64  `json:"total_value"`
	} `json:"data"`
}

type Fee struct {
	Display string `json:",omitempty"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Amount  int64  `json:"amount"`
}
