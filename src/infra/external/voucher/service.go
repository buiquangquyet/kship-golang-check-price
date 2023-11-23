package voucherext

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/dto"
	"check-price/src/infra/external"
	"context"
	"github.com/imroc/req/v3"
	"strconv"
	"time"
)

const (
	timeoutVoucher   = 5 * time.Second
	checkVoucherPath = "/vouchers/check"
)

type Service struct {
	*external.BaseClient
	client *req.Client
}

func NewService(base *external.BaseClient) *Service {
	cf := configs.Get().ExtService.Voucher
	cli := req.C().SetBaseURL(cf.Host).SetTimeout(timeoutVoucher)
	cli.SetCommonHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	base.SetTracer(cli)
	return &Service{
		BaseClient: base,
		client:     cli,
	}
}

func (g *Service) api(ctx context.Context) *req.Request {
	return g.client.R().SetContext(ctx)
}

func (g *Service) CheckVoucher(ctx context.Context, code string, retailerId, clientId int64) (*dto.Voucher, *common.Error) {
	var output checkVoucherOutput
	resp, err := g.api(ctx).
		SetFormData(map[string]string{
			"retailer_id": strconv.FormatInt(retailerId, 10),
			"client_id":   strconv.FormatInt(clientId, 10),
			"code":        code,
		}).
		SetSuccessResult(&output).
		SetErrorResult(&output).
		Post(checkVoucherPath)
	if err != nil {
		return nil, common.ErrSystemError(ctx, err.Error()).SetSource(common.SourceGHTKService)
	}

	if resp.IsErrorState() {
		log.Debug(ctx, "Call GetCompany MISA failed with body: %+v", output)
		//Todo consider error code
		//detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		//return nil, g.handleError(ctx, resp.StatusCode, &output).SetSource(common.SourceGHTKService).SetDetail(detail)
	}
	return output.ToDTO(), nil
}
