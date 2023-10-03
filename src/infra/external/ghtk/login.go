package ghtk

import "check-price/src/core/domain"

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newLoginInput(shop *domain.Shop) *LoginInput {
	return &LoginInput{
		Username: shop.GHTKUsername,
		Password: shop.GHTKPassword,
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
