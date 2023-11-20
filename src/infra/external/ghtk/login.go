package ghtkext

import "check-price/src/core/domain"

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newLoginInput(shop *domain.Shop) *LoginInput {
	return &LoginInput{
		Email:    shop.GHTKUsername,
		Password: shop.GHTKPass,
	}
}

type LoginOutput struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Code  string `json:"code"`
		Token string `json:"token"`
	} `json:"data"`
}
