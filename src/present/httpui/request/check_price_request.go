package request

type GetPriceReRequest struct {
	ClientCode         string          `uri:"client-code" binding:"required"`
	ActiveKShip        bool            `json:"ACTIVE_KSHIP"`
	SenderLocationId   string          `json:"SENDER_LOCATION_ID"`
	SenderWardId       string          `json:"SENDER_WARD_ID"`
	SenderAddress      string          `json:"SENDER_ADDRESS"`
	ReceiverLocationId int64           `json:"RECEIVER_LOCATION_ID"`
	ReceiverWardId     string          `json:"RECEIVER_WARD_ID"`
	ReceiverAddress    string          `json:"RECEIVER_ADDRESS"`
	ProductWidth       string          `json:"PRODUCT_WIDTH"`
	ProductHeight      string          `json:"PRODUCT_HEIGHT"`
	ProductLength      string          `json:"PRODUCT_LENGTH"`
	ProductWeight      string          `json:"PRODUCT_WEIGHT"`
	ProductQuantity    string          `json:"PRODUCT_QUANTITY"`
	ProductPrice       string          `json:"PRODUCT_PRICE"`
	MoneyCollection    string          `json:"MONEY_COLLECTION"`
	Services           []*Service      `json:"SERVICES"`
	ExtraService       []*ExtraService `json:"SERVICE_EXTRA"`
	VersionLocation    int             `json:"version_location"`
}

type Service struct {
	Code          string `json:"CODE"`
	OldTotalPrice string `json:"OLD_TOTAL_PRICE"`
}

type ExtraService struct {
	Code     string `json:"Code"`
	Value    string `json:"Value"`
	ViewType int    `json:"ViewType"`
	Name     string `json:"Name"`
}
