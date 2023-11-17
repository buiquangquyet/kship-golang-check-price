package aieliminating

type LoginOutput struct {
	Data struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

type OutputError struct {
	Message    string `json:"message"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}
