package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
	"check-price/src/helpers"
	"context"
)

func (p *PriceService) validate(ctx context.Context, v *dto.ValidatePrice) *common.Error {
	ierr := p.validateShop(ctx, v.Shop, v.ClientCode)
	if ierr != nil {
		return ierr
	}
	ierr = p.validateClient(ctx, v.ClientCode, v.RetailerId)
	if ierr != nil {
		return ierr
	}
	ierr = p.validateLocation(ctx, v.ClientCode, v.ValidateLocation)
	if ierr != nil {
		return ierr
	}
	ierr = p.validateService(ctx, v.ClientCode, v.Services)
	if ierr != nil {
		return ierr
	}
	ierr = p.validateExtraService(ctx, v.ClientCode, v.RetailerId, v.ExtraServices)
	if ierr != nil {
		return ierr
	}
	return nil
}

// done
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
		if shop.GHTKUsername == "" || shop.GHTKPassword == "" {
			return ierr.SetCode(3012)
		}
	case constant.JTFWDeliveryCode:
		if shop.JtCustomerId == "" {
			return ierr.SetCode(3005)
		}
	case constant.GHNFWDeliveryCode:
		if shop.GHNGWShopId == "" || shop.GHNFWPhone == "" {
			return ierr.SetCode(3005)
		}
	case constant.BESTFWDeliveryCode:
		if shop.UsernameBestFw == "" || shop.PasswordBestFw == "" {
			return ierr.SetCode(3005)
		}
	}
	return nil
}

// done
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
	clientUnAllowedShop, err := p.settingShopRepo.GetByRetailerId(ctx, retailerId)
	if err != nil {
		log.IErr(ctx, err)
		return err
	}
	if helpers.InArray(clientUnAllowedShop, client.Id) {
		return ierr.SetCode(3004)
	}

	if client.OnBoardingStatus == constant.OnboardingDisable {
		clientSettingShop, err := p.settingShopRepo.GetEnableShopByRetailerId(ctx, retailerId)
		if err != nil {
			log.IErr(ctx, err)
			return err
		}
		if !helpers.InArray(clientSettingShop, client.Id) {
			return ierr.SetCode(3004)
		}
	}
	return nil
}

func (p *PriceService) validateLocation(ctx context.Context, clientCode string, req *dto.ValidateLocation) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if req.VersionLocation == constant.VersionLocation2 {
		_, ierr = p.districtRepo.GetByKmsId(ctx, req.ReceiverLocationId)
	} else {
		_, ierr = p.districtRepo.GetByKvId(ctx, req.ReceiverLocationId)
	}
	if ierr != nil {
		return ierr.SetCode(4003)
	}
	if clientCode == constant.GHTKDeliveryCode && req.SenderWardId == 0 {
		return ierr.SetCode(4004)
	}
	if req.SenderWardId != 0 {
		//Todo duplicate code
		if req.VersionLocation == constant.VersionLocation2 {
			_, ierr = p.wardRepo.GetByKmsId(ctx, req.SenderWardId)
		} else {
			_, ierr = p.wardRepo.GetByKvId(ctx, req.SenderWardId)
		}
		if ierr != nil {
			return ierr.SetCode(4004)
		}
	}
	ierr = common.ErrBadRequest(ctx)
	if !helpers.InArray(constant.ReceiverWardIdClientCode, clientCode) && req.ReceiverWardId != 0 {
		return ierr.SetCode(4006)
	}
	//validate width, height, ...
	return nil
}

func (p *PriceService) validateService(ctx context.Context, clientCode string, services []*dto.ValidateService) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if services == nil {
		return ierr.SetCode(4001)
	}
	serviceEnable, ierr := p.serviceRepo.GetServicesPluckCodeByClientCode(ctx, clientCode)
	if ierr != nil {
		log.IErr(ctx, ierr)
		return ierr
	}
	for _, service := range services {
		serviceCode := service.Code
		if helpers.InArray(serviceEnable, serviceCode) {
			return common.ErrBadRequest(ctx).SetCode(4009)
		}
	}
	return nil
}

func (p *PriceService) validateExtraService(ctx context.Context, clientCode string, retailerId int64, extraServices []*dto.ValidateExtraService) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	extraServiceRequestCodes := make([]string, len(extraServices))
	for i, service := range extraServices {
		extraServiceRequestCodes[i] = service.Code
	}
	if !helpers.InArray(extraServiceRequestCodes, constant.ServiceExtraCodePayment) {
		return ierr.SetCode(4013)
	}
	servicesExtraByClientCodes, ierr := p.serviceRepo.GetByServiceCode(ctx, clientCode)
	if ierr != nil {
		log.IErr(ctx, ierr)
		return ierr
	}
	clientExtraServicesAllow := make([]string, 0)
	clientExtraServicesPaymentByAllow := make([]string, 0)
	for _, servicesExtraByClientCode := range servicesExtraByClientCodes {
		if servicesExtraByClientCode.Code == constant.ServiceExtraCodePayment {
			serviceExtraValue := constant.PaymentByFrom
			if servicesExtraByClientCode.Value != "" {
				serviceExtraValue = servicesExtraByClientCode.Value
			}
			clientExtraServicesPaymentByAllow = append(clientExtraServicesPaymentByAllow, serviceExtraValue)
		}
		if servicesExtraByClientCode.OnBoardingStatus == constant.OnboardingEnable {
			clientExtraServicesAllow = append(clientExtraServicesAllow, servicesExtraByClientCode.Code)
			continue
		}
		if servicesExtraByClientCode.OnBoardingStatus == constant.OnboardingDisable {
			serviceExtraEnableShop, ierr := p.settingShopRepo.GetServiceExtraEnableShop(ctx, retailerId)
			if ierr != nil {
				log.IErr(ctx, ierr)
				return ierr
			}
			if serviceExtraEnableShop {
				clientExtraServicesAllow = append(clientExtraServicesAllow, servicesExtraByClientCode.Code)
			}
		}
	}
	for _, extraService := range extraServices {
		extraServiceCode := extraService.Code
		if extraServiceCode == constant.ServiceExtraPickShift && !helpers.InArray(constant.PickShiftExtraClientCode, extraServiceCode) {
			return ierr.SetCode(4545)
		}
		if extraServiceCode == constant.ServiceExtraXmg && clientCode != constant.VTPFWDeliveryCode {
			return ierr.SetCode(4546)
		}
		if !helpers.InArray(clientExtraServicesAllow, extraServiceCode) {
			switch extraServiceCode {
			case constant.ServiceExtraBulky:
				return ierr.SetCode(4044)
			case constant.ServiceExtraCODST:
				return ierr.SetCode(4010)
			case constant.ServiceExtraConn:
				return ierr.SetCode(4011)
			case constant.ServiceExtraPartSign:
				return ierr.SetCode(4016)
			case constant.ServiceExtraThermalBag:
				return ierr.SetCode(4017)
			case constant.ServiceExtraGbh:
				return ierr.SetCode(4018)
			case constant.ServiceExtraShipperNote:
				return ierr.SetCode(4019)
			case constant.ServiceExtraGnG:
				return ierr.SetCode(4020)
			case constant.ServiceExtraRoundTrip:
				return ierr.SetCode(4021)
			case constant.ServiceExtraPrepaid:
				return ierr.SetCode(4022)
			case constant.ServiceExtraBaoPhat:
				return ierr.SetCode(4023)
			case constant.ServiceExtraPtt:
				return ierr.SetCode(4024)
			case constant.ServiceExtraDk:
				return ierr.SetCode(4025)
			case constant.ServiceExtraGHTK1:
				return ierr.SetCode(4402)
			case constant.ServiceExtraGHTK7:
				return ierr.SetCode(4403)
			case constant.ServiceExtraGHTKXFAST:
				return ierr.SetCode(4404)
			case constant.ServiceExtraGHTK2:
				return ierr.SetCode(4412)
			case constant.ServiceExtraGHTK3:
				return ierr.SetCode(4411)
			case constant.ServiceExtraGHTK4:
				return ierr.SetCode(4410)
			case constant.ServiceExtraGHTK5:
				return ierr.SetCode(4409)
			}
		}
		if extraServiceCode == constant.ServiceExtraCodePayment {
			serviceExtraValue := constant.PaymentByFrom
			if extraService.Value != "" {
				serviceExtraValue = extraService.Value
			}
			notAllow := !helpers.InArray(clientExtraServicesAllow, serviceExtraValue)
			if serviceExtraValue == constant.PaymentByFrom && notAllow {
				return ierr.SetCode(4001)
			}
			if serviceExtraValue == constant.PaymentByTo && notAllow {
				return ierr.SetCode(4002)
			}
		}
	}

	return nil
}
