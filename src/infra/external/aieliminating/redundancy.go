package aieliminating

type RedundancyOutput struct {
	Data struct {
		AddressNew string `json:"address_new"`
		Flag       string `json:"flag"`
	} `json:"data"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}
