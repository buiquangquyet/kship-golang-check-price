package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
	"check-price/src/core/enums"
	"check-price/src/helpers"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ExtraService struct {
	configCofT0Repo domain.ConfigCodT0Repo
	serviceRepo     domain.ServiceRepo
	settingShopRepo domain.SettingShopRepo
	settingRepo     domain.SettingRepo
	shopCodT0Repo   domain.ShopCodT0Repo
}

func NewExtraService(
	configCofT0Repo domain.ConfigCodT0Repo,
	serviceRepo domain.ServiceRepo,
	settingShopRepo domain.SettingShopRepo,
	settingRepo domain.SettingRepo,
	shopCodT0Repo domain.ShopCodT0Repo,
) *ExtraService {
	return &ExtraService{
		configCofT0Repo: configCofT0Repo,
		serviceRepo:     serviceRepo,
		settingShopRepo: settingShopRepo,
		settingRepo:     settingRepo,
		shopCodT0Repo:   shopCodT0Repo,
	}
}

func (c *ExtraService) handlePriceSpecialService(ctx context.Context, price *domain.Price, addInfoDto *dto.AddInfoDto) *common.Error {
	extraServiceCode := make([]string, len(addInfoDto.ExtraService))
	//payer := ""
	for i, service := range addInfoDto.ExtraService {
		if service.Code == "PAYMENT_BY" {
			//payer = service.Code
		}
		extraServiceCode[i] = service.Code
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODST) && c.checkServiceExtraIsPossible(ctx, addInfoDto, constant.ServiceExtraCODST) {
		c.addCodStPrice(price, addInfoDto.Shop, addInfoDto.Cod)
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODT0) {
		err := c.addCodT0Price(ctx, price, addInfoDto)
		if err != nil {
			return err
		}
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraConn) && c.checkServiceExtraIsPossible(ctx, addInfoDto, constant.ServiceExtraConn) {

	}
	return nil
}

func (c *ExtraService) addCodStPrice(price *domain.Price, shop *domain.Shop, cod int64) {
	var codStFee int64
	isShopType := shop.Type == constant.ShopVip
	for i := 0; i < constant.MaxLevel; i++ {
		if constant.CodLevelMin[i] <= cod && cod <= constant.CodLevelMax[i] {
			if isShopType {
				codStFee = constant.PriceVip[i]
			} else {
				codStFee = constant.PriceNormal[i]
			}
		}
	}
	price.CodstFee = codStFee
}

func (c *ExtraService) addCodT0Price(ctx context.Context, price *domain.Price, addInfoDto *dto.AddInfoDto) *common.Error {
	configCodT0s, isValid, err := c.validateCODT0(ctx, price, addInfoDto)
	if err != nil {
		return err
	}
	status := true
	msg := "success"
	var dataFee float64 = 0
	if !isValid {
		status = false
		msg = "Số tiền thu hộ hoặc dịch vụ không được áp dụng Đối soát nhanh."
	} else {
		dataFee, err = c.calculator(ctx, configCodT0s[0], addInfoDto)
		if err != nil {
			return err
		}
	}
	price.SetCodT0Info(status, msg, dataFee)
	return nil
}

func (c *ExtraService) calculator(ctx context.Context, configCodT0 *domain.ConfigCodT0, addInfoDto *dto.AddInfoDto) (float64, *common.Error) {
	isTrial, err := c.isTrial(ctx, addInfoDto.Shop)
	if err != nil {
		return 0, err
	}
	if isTrial {
		return 0, nil
	}
	value := configCodT0.Value
	switch configCodT0.Type {
	case constant.TypeFixed:
		return value, nil
	case constant.TypeByCodPercent:
		return float64(addInfoDto.Cod) * value / 100, nil
	}
	return 0, nil
}

func (c *ExtraService) isTrial(ctx context.Context, shop *domain.Shop) (bool, *common.Error) {
	settingFreeDay, ierr := c.settingRepo.GetByName(ctx, "free_trial_cod_t0")
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return false, ierr
	}
	feeDay, err := strconv.Atoi(settingFreeDay.Value)
	if err != nil {
		log.Error(ctx, err.Error())
		return false, common.ErrSystemError(ctx, err.Error())
	}
	shopCodT0, ierr := c.shopCodT0Repo.GetByShopId(ctx, shop.Id)
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return false, ierr
	}
	startDate := shopCodT0.TimeStart.Add(24 * time.Hour * time.Duration(feeDay))
	return time.Now().Before(startDate), nil
}

func (c *ExtraService) validateCODT0(ctx context.Context, price *domain.Price, addInfoDto *dto.AddInfoDto) ([]*domain.ConfigCodT0, bool, *common.Error) {
	configCodT0s, err := c.configCofT0Repo.GetByCodAndClientId(ctx, addInfoDto.Cod, addInfoDto.Client.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, false, err
	}
	if len(configCodT0s) == 0 {
		return configCodT0s, false, nil
	}

	codT0Service, isValidate, err := c.validateService(ctx, addInfoDto.RetailerId, addInfoDto.Client.Code, price)
	if err != nil {
		return configCodT0s, false, err
	}
	if !isValidate {
		return configCodT0s, false, nil
	}

	shopBlackList, err := c.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraDisableShop, addInfoDto.RetailerId)
	if err != nil {
		log.Error(ctx, err.Error())
		return configCodT0s, false, err
	}
	if len(shopBlackList) != 0 {
		return configCodT0s, false, nil
	}
	isUseCodT0, err := c.settingShopRepo.GetByRetailerIdAndModelId(ctx, enums.ModelTypeServiceExtraSettingUser, addInfoDto.RetailerId, codT0Service.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		return configCodT0s, false, nil
	}
	if len(isUseCodT0) == 0 {
		return configCodT0s, false, nil
	}
	return configCodT0s, true, nil
}

func (c *ExtraService) validateService(ctx context.Context, retailerId int64, clientCode string, price *domain.Price) (*domain.Service, bool, *common.Error) {
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

func (c *ExtraService) checkServiceExtraIsPossible(ctx context.Context, addInfoDto *dto.AddInfoDto, extraServiceCode string) bool {
	extraService, err := c.serviceRepo.GetByCode(ctx, extraServiceCode)
	if helpers.IsInternalError(err) {
		log.Error(ctx, err.Error())
		return false
	}
	if err != nil {
		return false
	}
	if extraService.ClientsPossible == "" && strings.Contains(extraService.ClientsPossible, addInfoDto.Client.Code) {
		if extraService.OnBoardingStatus == constant.StatusEnableServiceExtra {
			return true
		}
		serviceExtraEnableShop, err := c.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraSettingShop, addInfoDto.RetailerId)
		if err != nil {
			log.Error(ctx, err.Error())
			return false
		}
		for _, service := range serviceExtraEnableShop {
			if extraService.Id == service.ModelId {
				return true
			}
		}
	}
	return true
}
