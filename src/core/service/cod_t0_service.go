package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/helpers"
	"context"
	"fmt"
	"strconv"
	"strings"
)

type CodT0Service struct {
	configCofT0Repo domain.ConfigCodT0Repo
	serviceRepo     domain.ServiceRepo
	settingShopRepo domain.SettingShopRepo
}

func NewCodT0Service(
	configCofT0Repo domain.ConfigCodT0Repo,
	serviceRepo domain.ServiceRepo,
	settingShopRepo domain.SettingShopRepo,
) *CodT0Service {
	return &CodT0Service{
		configCofT0Repo: configCofT0Repo,
		serviceRepo:     serviceRepo,
		settingShopRepo: settingShopRepo,
	}
}

//func (c *CodT0Service) addCodT0Price(ctx context.Context, price *domain.Price) *common.Error {
//	isValid, err := c.validateCODT0(ctx, price)
//	if err != nil {
//		return err
//	}
//	status := true
//	msg := "success"
//	if !isValid {
//		status = false
//		msg = "Số tiền thu hộ không được áp dụng Đối soát nhanh."
//	}
//	price.
//}

func (c *CodT0Service) validateCODT0(ctx context.Context, cod int64, clientId int64, retailerId int64, clientCode string, price *domain.Price) (bool, *common.Error) {
	configCodT0s, err := c.configCofT0Repo.GetByCodAndClientId(ctx, cod, clientId)
	if err != nil {
		log.Error(ctx, err.Error())
		return false, err
	}
	if len(configCodT0s) == 0 {
		return false, nil
	}

	codT0Service, isValidate, err := c.validateService(ctx, retailerId, clientCode, price)
	if err != nil {
		return false, err
	}
	if !isValidate {
		return false, nil
	}

	shopBlackList, err := c.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraDisableShop, retailerId)
	if err != nil {
		log.Error(ctx, err.Error())
		return false, err
	}
	if len(shopBlackList) != 0 {
		return false, nil
	}
	isUseCodT0, err := c.settingShopRepo.GetByRetailerIdAndModelId(ctx, enums.ModelTypeServiceExtraSettingUser, retailerId, codT0Service.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		return false, nil
	}
	if len(isUseCodT0) == 0 {
		return false, nil
	}
	return true, nil
}

func (c *CodT0Service) validateService(ctx context.Context, retailerId int64, clientCode string, price *domain.Price) (*domain.Service, bool, *common.Error) {
	codT0Service, err := c.serviceRepo.GetByCode(ctx, constant.ServiceExtraCODT0)
	if helpers.IsInternalError(err) {
		log.Error(ctx, err.Error())
		return nil, false, err
	}
	if err != nil {
		return nil, false, nil
	}
	if codT0Service.Status == constant.StatusDisableServiceExtra {
		return nil, false, nil
	}
	if codT0Service.OnBoardingStatus == constant.OnboardingDisable {
		shopSettings, err := c.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraSettingShop, retailerId)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, false, err
		}
		fmt.Println(shopSettings)
	}
	idPrice := strconv.FormatInt(price.Id, 10)
	if strings.Contains(codT0Service.ClientsPossible, clientCode) || strings.Contains(codT0Service.ClientsPossible, idPrice) {
		return nil, false, nil
	}
	return codT0Service, true, nil
}
