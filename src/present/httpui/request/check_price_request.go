package request

type GetPriceReRequest struct {
	ClientCode         string          `uri:"client" binding:"required"`
	ActiveKShip        bool            `json:"ACTIVE_KSHIP"`
	SenderLocationId   int             `json:"SENDER_LOCATION_ID"`
	SenderWardId       int             `json:"SENDER_WARD_ID"`
	SenderAddress      string          `json:"SENDER_ADDRESS"`
	ReceiverLocationId int64           `json:"RECEIVER_LOCATION_ID"`
	ReceiverWardId     int             `json:"RECEIVER_WARD_ID"`
	ReceiverAddress    string          `json:"RECEIVER_ADDRESS"`
	ProductWidth       int             `json:"PRODUCT_WIDTH"`
	ProductHeight      int             `json:"PRODUCT_HEIGHT"`
	ProductLength      int             `json:"PRODUCT_LENGTH"`
	ProductWeight      int             `json:"PRODUCT_WEIGHT"`
	ProductQuantity    int             `json:"PRODUCT_QUANTITY"`
	ProductPrice       string          `json:"PRODUCT_PRICE"`
	MoneyCollection    int64           `json:"MONEY_COLLECTION"`
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
	ViewType string `json:"ViewType"`
	Name     string `json:"Name"`
}
