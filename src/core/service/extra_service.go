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
	shopLevelRepo   domain.ShopLevelRepo
}

func NewExtraService(
	configCofT0Repo domain.ConfigCodT0Repo,
	serviceRepo domain.ServiceRepo,
	settingShopRepo domain.SettingShopRepo,
	settingRepo domain.SettingRepo,
	shopCodT0Repo domain.ShopCodT0Repo,
	shopLevelRepo domain.ShopLevelRepo,
) *ExtraService {
	return &ExtraService{
		configCofT0Repo: configCofT0Repo,
		serviceRepo:     serviceRepo,
		settingShopRepo: settingShopRepo,
		settingRepo:     settingRepo,
		shopCodT0Repo:   shopCodT0Repo,
		shopLevelRepo:   shopLevelRepo,
	}
}

func (s *ExtraService) handlePriceSpecialService(ctx context.Context, price *domain.Price, addInfoDto *dto.AddInfoDto) *common.Error {
	extraServiceCode := make([]string, len(addInfoDto.ExtraService))
	for i, service := range addInfoDto.ExtraService {
		extraServiceCode[i] = service.Code
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODST) && s.checkServiceExtraIsPossible(ctx, addInfoDto, constant.ServiceExtraCODST) {
		s.addCodStPrice(price, addInfoDto.Shop, addInfoDto.Cod)
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODT0) {
		err := s.addCodT0Price(ctx, price, addInfoDto)
		if err != nil {
			return err
		}
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraConn) && s.checkServiceExtraIsPossible(ctx, addInfoDto, constant.ServiceExtraConn) {
		err := s.addCodConnPrice(ctx, price, addInfoDto.Shop)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ExtraService) addCodStPrice(price *domain.Price, shop *domain.Shop, cod int64) {
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
	price.SetCodStFee(codStFee)
}

func (s *ExtraService) addCodConnPrice(ctx context.Context, price *domain.Price, shop *domain.Shop) *common.Error {
	var shopLevelId int64 = 1
	var err error
	if shop.ShopLevel != "" {
		shopLevelId, err = strconv.ParseInt(shop.ShopLevel, 10, 64)
		if err != nil {
			log.Error(ctx, err.Error())
			return common.ErrSystemError(ctx, err.Error())
		}
	}
	shopLevel, ierr := s.shopLevelRepo.GetById(ctx, shopLevelId)
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return ierr
	}
	price.SetConnFee(shopLevel.GhnMarkup)
	return nil
}
func (s *ExtraService) addCodT0Price(ctx context.Context, price *domain.Price, addInfoDto *dto.AddInfoDto) *common.Error {
	configCodT0s, isValid, err := s.validateCODT0(ctx, price, addInfoDto)
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
		dataFee, err = s.calculator(ctx, configCodT0s[0], addInfoDto)
		if err != nil {
			return err
		}
	}
	price.SetCodT0Info(status, msg, dataFee)
	return nil
}

func (s *ExtraService) calculator(ctx context.Context, configCodT0 *domain.ConfigCodT0, addInfoDto *dto.AddInfoDto) (float64, *common.Error) {
	isTrial, err := s.isTrial(ctx, addInfoDto.Shop)
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

func (s *ExtraService) isTrial(ctx context.Context, shop *domain.Shop) (bool, *common.Error) {
	settingFreeDay, ierr := s.settingRepo.GetByName(ctx, "free_trial_cod_t0")
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return false, ierr
	}
	feeDay, err := strconv.Atoi(settingFreeDay.Value)
	if err != nil {
		log.Error(ctx, err.Error())
		return false, common.ErrSystemError(ctx, err.Error())
	}
	shopCodT0, ierr := s.shopCodT0Repo.GetByShopId(ctx, shop.Id)
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return false, ierr
	}
	startDate := shopCodT0.TimeStart.Add(24 * time.Hour * time.Duration(feeDay))
	return time.Now().Before(startDate), nil
}

func (s *ExtraService) validateCODT0(ctx context.Context, price *domain.Price, addInfoDto *dto.AddInfoDto) ([]*domain.ConfigCodT0, bool, *common.Error) {
	configCodT0s, err := s.configCofT0Repo.GetByCodAndClientId(ctx, addInfoDto.Cod, addInfoDto.Client.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, false, err
	}
	if len(configCodT0s) == 0 {
		return configCodT0s, false, nil
	}

	codT0Service, isValidate, err := s.validateService(ctx, addInfoDto.RetailerId, addInfoDto.Client.Code, price)
	if err != nil {
		return configCodT0s, false, err
	}
	if !isValidate {
		return configCodT0s, false, nil
	}

	shopBlackList, err := s.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraDisableShop, addInfoDto.RetailerId)
	if err != nil {
		log.Error(ctx, err.Error())
		return configCodT0s, false, err
	}
	if len(shopBlackList) != 0 {
		return configCodT0s, false, nil
	}
	isUseCodT0, err := s.settingShopRepo.GetByRetailerIdAndModelId(ctx, enums.ModelTypeServiceExtraSettingUser, addInfoDto.RetailerId, codT0Service.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		return configCodT0s, false, nil
	}
	if len(isUseCodT0) == 0 {
		return configCodT0s, false, nil
	}
	return configCodT0s, true, nil
}

func (s *ExtraService) validateService(ctx context.Context, retailerId int64, clientCode string, price *domain.Price) (*domain.Service, bool, *common.Error) {
	codT0Service, err := s.serviceRepo.GetByCode(ctx, constant.ServiceExtraCODT0)
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
		_, err := s.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraSettingShop, retailerId)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, false, err
		}
	}
	idPrice := strconv.FormatInt(price.Id, 10)
	if strings.Contains(codT0Service.ClientsPossible, clientCode) || strings.Contains(codT0Service.ClientsPossible, idPrice) {
		return nil, false, nil
	}
	return codT0Service, true, nil
}

func (s *ExtraService) checkServiceExtraIsPossible(ctx context.Context, addInfoDto *dto.AddInfoDto, extraServiceCode string) bool {
	extraService, err := s.serviceRepo.GetByCode(ctx, extraServiceCode)
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
		serviceExtraEnableShop, err := s.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraSettingShop, addInfoDto.RetailerId)
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
