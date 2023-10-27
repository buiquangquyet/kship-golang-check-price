package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type Shop struct {
	Id          int64  `json:"id"`
	CityId      int64  `json:"city_id"`
	DistrictId  int64  `json:"district_id"`
	WardId      int64  `json:"ward_id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Code        string `json:"code"`
	ShopLevel   string `json:"shop_level"`
	Phone       string `json:"phone"`
	Sex         int    `json:"sex"`
	Address     string `json:"address"`
	Token       string `json:"token"`
	BankId      int    `json:"bank_id"`
	BankCode    string `json:"bank_code"`
	BankName    string `json:"bank_name"`
	BankBranch  string `json:"bank_branch"`
	BankAccount string `json:"bank_account"`
	BankNumber  string `json:"bank_number"`

	GhnShopId       string `json:"ghn_shop_id"`
	GhnContractId   int    `json:"ghn_contract_id"`
	VtpUsername     string `json:"vtp_username"`
	VtpPassword     string `json:"vtp_password"`
	VnpCmsCode      string `json:"vnp_cms_code"`
	VnpCustomerCode string `json:"vnp_customer_code"`
	GHTKUsername    string `json:"ghtk_username"`
	GHTKPassword    string `json:"ghtk_password"`
	GHTKPass        string `json:"ghtk_pass"`
	JtCustomerId    string `json:"jt_customer_id"`
	GHNGWShopId     string `json:"ghnfw_shop_id"`
	GHNFWPhone      string `json:"ghnfw_phone"`
	UsernameBestFw  string `json:"username_bestfw"`
	PasswordBestFw  string `json:"password_bestfw"`

	Type int `json:"type"`
	//update sau
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type ShopRepo interface {
	GetByRetailerId(ctx context.Context, retailerId int64) (*Shop, *common.Error)
	GetByCode(ctx context.Context, shopCode string) (*Shop, *common.Error)
}

func (Shop) TableName() string {
	return "shops"
}
