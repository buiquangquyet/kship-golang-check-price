package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
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

func (s *VoucherService) checkVoucher(ctx context.Context, price *domain.Price, addInfoDto *dto.AddInfoDto) *common.Error {
	if addInfoDto.Coupon == "" {
		return nil
	}
	voucher, err := s.voucherExtService.CheckVoucher(ctx, addInfoDto.Coupon, addInfoDto.RetailerId, addInfoDto.Client.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		return err
	}
	callTo := 0
	//bo qua client_error_code
	switch voucher.StatusCode {
	case constant.VoucherExist:
		if voucher.Type == constant.TypeVoucherKv && !(addInfoDto.Payer == "NGUOINHAN") {
			callTo = constant.UseKv
		} else if voucher.Type == constant.TypeVoucherDelivery {
			if helpers.InArray(constant.ClientAllowUsePromotion, addInfoDto.Client.Code) {
				callTo = constant.UseDelivery
			}
		}
	case constant.VoucherNotExist:
		if helpers.InArray(constant.ClientAllowUsePromotion, addInfoDto.Client.Code) {
			callTo = constant.UseDelivery
		}
	case constant.VoucherError:
		return nil
	}
	if callTo != constant.UseKv {
		return nil
	}
	discountVoucher := voucher.DiscountValue
	if discountVoucher > price.TotalPrice {
		discountVoucher = price.TotalPrice
	}
	price.SetCouponInfo(discountVoucher, price.TotalPrice)
	return nil
}
