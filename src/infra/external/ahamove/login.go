package ahamoveext

type LoginOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type OutputErr struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
