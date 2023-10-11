package request

import (
	"check-price/src/common"
	"context"
)

type GetPriceReRequest struct {
	ClientCode  string `uri:"client" binding:"required"`
	ActiveKShip bool   `json:"ACTIVE_KSHIP"`
	*Sender
	*Receiver
	*Product
	MoneyCollection int64           `json:"MONEY_COLLECTION"`
	Services        []*Service      `json:"SERVICES"`
	ExtraService    []*ExtraService `json:"SERVICE_EXTRA"`

	RetailerId      int64 `json:"retailer_id"`
	VersionLocation int   `json:"version_location"`
}

type Sender struct {
	SenderLocationId int64  `json:"SENDER_LOCATION_ID"`
	SenderWardId     int64  `json:"SENDER_WARD_ID"`
	SenderAddress    string `json:"SENDER_ADDRESS"`
}

type Receiver struct {
	ReceiverLocationId int64  `json:"RECEIVER_LOCATION_ID"`
	ReceiverWardId     int64  `json:"RECEIVER_WARD_ID"`
	ReceiverAddress    string `json:"RECEIVER_ADDRESS"`
}

type Product struct {
	ProductWidth    int   `json:"PRODUCT_WIDTH"`
	ProductHeight   int   `json:"PRODUCT_HEIGHT"`
	ProductLength   int   `json:"PRODUCT_LENGTH"`
	ProductWeight   int64 `json:"PRODUCT_WEIGHT"`
	ProductQuantity int   `json:"PRODUCT_QUANTITY"`
	ProductPrice    int64 `json:"PRODUCT_PRICE"`
}

type Service struct {
	Code          string `json:"CODE"`
	OldTotalPrice string `json:"OLD_TOTAL_PRICE"`
}

type ExtraService struct {
	Code     string `json:"Code"`
	Value    string `json:"Value"`
	ViewType string `json:"ViewType"`
	Name     string `json:"Name"`
}

func (g *GetPriceReRequest) validate(ctx context.Context) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if g.ReceiverLocationId == 0 {
		return ierr.SetCode(4005)
	}
	if g.ProductWidth == 0 || g.ProductHeight == 0 || g.ProductLength == 0 || g.ProductWeight == 0 {
		return ierr.SetCode(4007)
	}
	if g.SenderLocationId == 0 {
		return ierr.SetCode(4003)
	}
	if len(g.ExtraService) == 0 {
		return ierr.SetCode(4013)
	}
	for _, extraService := range g.ExtraService {
		if extraService.Code == "" {
			return ierr.SetCode(4012)
		}
	}
	return nil
}
