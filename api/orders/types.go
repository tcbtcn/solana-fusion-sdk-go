package orders

// OrderStatus represents the status of an order
type OrderStatus int

const (
	OrderStatusInProgress OrderStatus = 0
	OrderStatusFilled     OrderStatus = 200
	OrderStatusCancelled  OrderStatus = 301
)

// OrderDTO represents an order
type OrderDTO struct {
	ID                          uint32              `json:"id"`
	SrcAmount                   string              `json:"srcAmount"`
	MinDstAmount                string              `json:"minDstAmount"`
	EstimatedDstAmount          string              `json:"estimatedDstAmount"`
	ExpirationTime              uint32              `json:"expirationTime"`
	Receiver                    string              `json:"receiver"`
	SrcAssetIsNative            bool                `json:"srcAssetIsNative"`
	DstAssetIsNative            bool                `json:"dstAssetIsNative"`
	Fee                         FeeDTO              `json:"fee"`
	DutchAuctionData            DutchAuctionDataDTO `json:"dutchAuctionData"`
	SrcMint                     string              `json:"srcMint"`
	DstMint                     string              `json:"dstMint"`
	CancellationAuctionDuration uint32              `json:"cancellationAuctionDuration"`
}

// FeeDTO represents fee information
type FeeDTO struct {
	ProtocolFee            uint16 `json:"protocolFee"`
	IntegratorFee          uint16 `json:"integratorFee"`
	SurplusPercentage      uint8  `json:"surplusPercentage"`
	ProtocolDstAta         string `json:"protocolDstAta"`
	IntegratorDstAta       string `json:"integratorDstAta"`
	MaxCancellationPremium string `json:"maxCancellationPremium"`
}

// DutchAuctionDataDTO represents Dutch auction data
type DutchAuctionDataDTO struct {
	StartTime           uint32                 `json:"startTime"`
	Duration            uint32                 `json:"duration"`
	InitialRateBump     uint16                 `json:"initialRateBump"`
	PointsAndTimeDeltas []PointAndTimeDeltaDTO `json:"pointsAndTimeDeltas"`
}

// PointAndTimeDeltaDTO represents a point and time delta
type PointAndTimeDeltaDTO struct {
	RateBump  uint16 `json:"rateBump"`
	TimeDelta uint16 `json:"timeDelta"`
}

// OrderInfoDTO represents order information
type OrderInfoDTO struct {
	OrderHash            string   `json:"orderHash"`
	TxSignature          string   `json:"txSignature"`
	Maker                string   `json:"maker"`
	Order                OrderDTO `json:"order"`
	RemainingMakerAmount string   `json:"remainingMakerAmount"`
}

// OrderCancellableByResolverInfoDTO represents order cancellable by resolver
type OrderCancellableByResolverInfoDTO struct {
	OrderHash   string   `json:"orderHash"`
	TxSignature string   `json:"txSignature"`
	Maker       string   `json:"maker"`
	Order       OrderDTO `json:"order"`
}

// FillDTO represents a fill
type FillDTO struct {
	TxSignature              string `json:"txSignature"`
	FilledMakerAmount        string `json:"filledMakerAmount"`
	FilledAuctionTakerAmount string `json:"filledAuctionTakerAmount"`
}

// OrderStatusDTO represents order status
type OrderStatusDTO struct {
	Maker                   string      `json:"maker"`
	OrderHash               string      `json:"orderHash"`
	Status                  OrderStatus `json:"status"`
	Order                   OrderDTO    `json:"order"`
	ApproximateTakingAmount string      `json:"approximateTakingAmount"`
	CancelTx                *string     `json:"cancelTx,omitempty"`
	ExpirationTime          uint32      `json:"expirationTime"`
	Fills                   []FillDTO   `json:"fills"`
	CreatedAt               uint32      `json:"createdAt"`
	SrcTokenPriceUsd        float64     `json:"srcTokenPriceUsd"`
	DstTokenPriceUsd        float64     `json:"dstTokenPriceUsd"`
	Cancelable              bool        `json:"cancelable"`
}
