package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/helpers"
	"check-price/src/present/httpui/request"
	"context"
)

type PriceService struct {
	shopRepo          domain.ShopRepo
	priceRepo         domain.PriceRepo
	clientRepo        domain.ClientRepo
	clientDisableShop domain.ClientDisableShopRepo
	clientSettingShop domain.ClientSettingShopRepo
	districtRepo      domain.DistrictRepo
}

func NewPriceService(shopRepo domain.ShopRepo) *PriceService {
	return &PriceService{
		shopRepo: shopRepo,
	}
}

func (p *PriceService) GetPrice(ctx context.Context, req *request.GetPriceReRequest) ([]*domain.Price, *common.Error) {
	shop, err := p.shopRepo.GetByRetailerId(ctx, "")
	if helpers.IsInternalError(err) {
		log.Error(ctx, err.Error())
		return nil, err
	}
	//tai sao dk nhu nay moi vao cache check
	if shop == nil || !req.ActiveKShip {
		// response từ cache
		prices, err := p.priceRepo.GetResponse(ctx, req.ClientCode, req.SenderWardId, req.ReceiverWardId)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, err
		}
		if prices != nil {
			return prices, nil
		}
		// get shop default, từ cache hay cả db phai hỏi lại
		shop, err = p.shopRepo.GetByCode(ctx, constant.ShopDefaultTrial)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, err
		}
	}

	return nil, nil
}

func (p *PriceService) validate(ctx context.Context, shop *domain.Shop, req *request.GetPriceReRequest) *common.Error {

	return nil
}

func (p *PriceService) validateShop(ctx context.Context, shop *domain.Shop, clientCode string) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	switch clientCode {
	case constant.VTPFWDeliveryCode:
		if shop.VtpUsername == "" || shop.VtpPassword == "" {
			return ierr.SetCode(3005)
		}
	case constant.VNPDeliveryCode:
		if shop.VnpCmsCode == "" {
			return ierr.SetCode(3008)
		}
	case constant.GHTKDeliveryCode:
		if shop.GhtkUsername == "" || shop.GhtkPassword == "" {
			return ierr.SetCode(3012)
		}
	case constant.JTFWDeliveryCode:
		//code cu bi duplicate
		if shop.JtCustomerId == "" {
			return ierr.SetCode(3005)
		}
	case constant.GHNFWDeliveryCode:
		if shop.GhnfwShopId == "" || shop.GhnfwPhone == "" {
			return ierr.SetCode(3005)
		}
	case constant.BESTFWDeliveryCode:
		if shop.UsernameBestfw == "" || shop.PasswordBestfw == "" {
			return ierr.SetCode(3005)
		}
	}
	return nil
}

func (p *PriceService) validateClient(ctx context.Context, clientCode string, retailerId int64) *common.Error {
	client, err := p.clientRepo.GetByCode(ctx, clientCode)
	if helpers.IsInternalError(err) {
		log.IErr(ctx, err)
		return err
	}
	ierr := common.ErrBadRequest(ctx)
	if client == nil {
		log.Warn(ctx, "client is null")
		return ierr.SetCode(3001)
	}
	if client.Status == constant.DisableStatus {
		return ierr.SetCode(3002)
	}
	clientUnAllowedShop, err := p.clientDisableShop.GetByRetailerId(ctx, retailerId)
	if err != nil {
		log.IErr(ctx, err)
		return err
	}
	if helpers.InArray(clientUnAllowedShop, client.Id) {
		return ierr.SetCode(3004)
	}

	if client.OnBoardingStatus == constant.OnboardingDisable {
		clientSettingShop, err := p.clientSettingShop.GetEnableShopByRetailerId(ctx, retailerId)
		if err != nil {
			log.IErr(ctx, err)
			return err
		}
		if helpers.InArray(clientSettingShop, client.Id) {
			return ierr.SetCode(3004)
		}
	}
	return nil
}

func (p *PriceService) validateLocation(ctx context.Context, req *request.GetPriceReRequest) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if req.SenderLocationId != "" {
		return ierr.SetCode(4003)
	}
	if req.VersionLocation == constant.VersionLocation2 {
		_, ierr = p.districtRepo.GetByKmsId(ctx, req.ReceiverLocationId)
	} else {
		_, ierr = p.districtRepo.GetByKvId(ctx, req.ReceiverLocationId)
	}
	if ierr != nil {
		return ierr.SetCode(4005)
	}
	ierr = common.ErrBadRequest(ctx)
	if !helpers.InArray(constant.ReceiverWardIdClientCode, req.ClientCode) && req.ReceiverWardId != "" {
		return ierr.SetCode(4006)
	}
	return nil
}
