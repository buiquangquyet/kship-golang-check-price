package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/dto"
	"check-price/src/core/enums"
	"check-price/src/helpers"
	"check-price/src/infra/external/voucher"
	"context"
)

type VoucherService struct {
	voucherExtService *voucher.VoucherExtService
}

func NewVoucherService(voucherExtService *voucher.VoucherExtService) *VoucherService {
	return &VoucherService{
		voucherExtService: voucherExtService,
	}
}

func (s *VoucherService) checkVoucher(ctx context.Context, addInfoDto *dto.AddInfoDto) (enums.TypeVoucherUse, int64, *common.Error) {
	if addInfoDto.Coupon == "" {
		return enums.TypeVoucherNotExist, 0, nil
	}
	voucherKv, err := s.voucherExtService.CheckVoucher(ctx, addInfoDto.Coupon, addInfoDto.RetailerId, addInfoDto.Client.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		//chay tiep
		return enums.TypeVoucherNotExist, 0, err
	}
	var callTo enums.TypeVoucherUse
	switch voucherKv.StatusCode {
	case constant.VoucherExist:
		if voucherKv.Type == constant.TypeVoucherKv && !(addInfoDto.Payer == constant.PaymentByTo) {
			callTo = enums.TypeVoucherUseKv
		}
		if voucherKv.Type == constant.TypeVoucherDelivery && helpers.InArray(constant.ClientAllowUsePromotion, addInfoDto.Client.Code) {
			callTo = enums.TypeVoucherUseDelivery
		}
	case constant.VoucherNotExist:
		if helpers.InArray(constant.ClientAllowUsePromotion, addInfoDto.Client.Code) {
			callTo = enums.TypeVoucherUseDelivery
		}
	case constant.VoucherError:
		//Todo log
		return enums.TypeVoucherNotExist, 0, nil
	}

	return callTo, voucherKv.DiscountValue, nil
}
