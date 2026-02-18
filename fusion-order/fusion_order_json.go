package fusionorder

import (
	"errors"
	"math/big"

	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
)

// FusionOrderJSON represents the JSON format of a FusionOrder
type FusionOrderJSON struct {
	ID                          uint32               `json:"id"`
	Receiver                    string               `json:"receiver"`
	CancellationAuctionDuration uint32               `json:"cancellationAuctionDuration"`
	SrcMint                     string               `json:"srcMint"`
	DstMint                     string               `json:"dstMint"`
	SrcAmount                   string               `json:"srcAmount"`
	MinDstAmount                string               `json:"minDstAmount"`
	EstimatedDstAmount          string               `json:"estimatedDstAmount"`
	ExpirationTime              uint32               `json:"expirationTime"`
	SrcAssetIsNative            bool                 `json:"srcAssetIsNative"`
	DstAssetIsNative            bool                 `json:"dstAssetIsNative"`
	Fee                         FeeJSON              `json:"fee"`
	DutchAuctionData            DutchAuctionDataJSON `json:"dutchAuctionData"`
}

type FeeJSON struct {
	ProtocolDstAta         *string `json:"protocolDstAta"`
	IntegratorDstAta       *string `json:"integratorDstAta"`
	ProtocolFee            uint16  `json:"protocolFee"`
	IntegratorFee          uint16  `json:"integratorFee"`
	SurplusPercentage      uint8   `json:"surplusPercentage"`
	MaxCancellationPremium string  `json:"maxCancellationPremium"`
}

type DutchAuctionDataJSON struct {
	StartTime           uint32                  `json:"startTime"`
	Duration            uint32                  `json:"duration"`
	InitialRateBump     uint16                  `json:"initialRateBump"`
	PointsAndTimeDeltas []PointAndTimeDeltaJSON `json:"pointsAndTimeDeltas"`
}

type PointAndTimeDeltaJSON struct {
	RateBump  uint16 `json:"rateBump"`
	TimeDelta uint16 `json:"timeDelta"`
}

// ToJSON converts FusionOrder to JSON format
func (f *FusionOrder) ToJSON() *FusionOrderJSON {
	order := f.Build()

	var protocolDstAta, integratorDstAta *string
	if f.orderConfig.Fees != nil {
		if f.orderConfig.Fees.ProtocolDstAta != nil {
			addr := f.orderConfig.Fees.ProtocolDstAta.ToString()
			protocolDstAta = &addr
		}
		if f.orderConfig.Fees.IntegratorDstAta != nil {
			addr := f.orderConfig.Fees.IntegratorDstAta.ToString()
			integratorDstAta = &addr
		}
	}

	pointsAndTimeDeltas := make([]PointAndTimeDeltaJSON, len(order.DutchAuctionData.PointsAndTimeDeltas))
	for i, point := range order.DutchAuctionData.PointsAndTimeDeltas {
		pointsAndTimeDeltas[i] = PointAndTimeDeltaJSON(point)
	}

	return &FusionOrderJSON{
		ID:                          order.ID,
		Receiver:                    f.orderConfig.Receiver.ToString(),
		CancellationAuctionDuration: order.CancellationAuctionDuration,
		SrcMint:                     f.orderConfig.SrcMint.ToString(),
		DstMint:                     f.orderConfig.DstMint.ToString(),
		SrcAmount:                   order.SrcAmount.String(),
		MinDstAmount:                order.MinDstAmount.String(),
		EstimatedDstAmount:          order.EstimatedDstAmount.String(),
		ExpirationTime:              order.ExpirationTime,
		SrcAssetIsNative:            order.SrcAssetIsNative,
		DstAssetIsNative:            order.DstAssetIsNative,
		Fee: FeeJSON{
			ProtocolDstAta:         protocolDstAta,
			IntegratorDstAta:       integratorDstAta,
			ProtocolFee:            order.Fee.ProtocolFee,
			IntegratorFee:          order.Fee.IntegratorFee,
			SurplusPercentage:      order.Fee.SurplusPercentage,
			MaxCancellationPremium: order.Fee.MaxCancellationPremium.String(),
		},
		DutchAuctionData: DutchAuctionDataJSON{
			StartTime:           order.DutchAuctionData.StartTime,
			Duration:            order.DutchAuctionData.Duration,
			InitialRateBump:     order.DutchAuctionData.InitialRateBump,
			PointsAndTimeDeltas: pointsAndTimeDeltas,
		},
	}
}

// FromJSON creates a FusionOrder from JSON
func FromJSON(jsonData *FusionOrderJSON) (*FusionOrder, error) {
	orderExpirationDelay := jsonData.ExpirationTime - jsonData.DutchAuctionData.Duration - jsonData.DutchAuctionData.StartTime

	// Parse addresses
	srcMint, err := domains.NewAddress(jsonData.SrcMint)
	if err != nil {
		return nil, err
	}
	dstMint, err := domains.NewAddress(jsonData.DstMint)
	if err != nil {
		return nil, err
	}
	receiver, err := domains.NewAddress(jsonData.Receiver)
	if err != nil {
		return nil, err
	}

	// Parse amounts
	srcAmount, ok := new(big.Int).SetString(jsonData.SrcAmount, 10)
	if !ok {
		return nil, errors.New("invalid srcAmount")
	}
	minDstAmount, ok := new(big.Int).SetString(jsonData.MinDstAmount, 10)
	if !ok {
		return nil, errors.New("invalid minDstAmount")
	}
	estimatedDstAmount, ok := new(big.Int).SetString(jsonData.EstimatedDstAmount, 10)
	if !ok {
		return nil, errors.New("invalid estimatedDstAmount")
	}

	// Parse fees
	var protocolDstAta, integratorDstAta *domains.Address
	if jsonData.Fee.ProtocolDstAta != nil {
		protocolDstAta, err = domains.NewAddress(*jsonData.Fee.ProtocolDstAta)
		if err != nil {
			return nil, err
		}
	}
	if jsonData.Fee.IntegratorDstAta != nil {
		integratorDstAta, err = domains.NewAddress(*jsonData.Fee.IntegratorDstAta)
		if err != nil {
			return nil, err
		}
	}

	protocolFee := domains.BpsFromFraction(float64(jsonData.Fee.ProtocolFee)/100000.0, Base1E5)
	integratorFee := domains.BpsFromFraction(float64(jsonData.Fee.IntegratorFee)/100000.0, Base1E5)
	surplusShare := domains.BpsFromFraction(float64(jsonData.Fee.SurplusPercentage)/100.0, Base1E2)

	fees, err := NewFeeConfig(protocolDstAta, integratorDstAta, protocolFee, integratorFee, surplusShare)
	if err != nil {
		return nil, err
	}

	// Parse resolver cancellation config
	var maxCancellationPremium *big.Int
	if jsonData.Fee.MaxCancellationPremium == "" {
		maxCancellationPremium = big.NewInt(0)
	} else {
		var ok bool
		maxCancellationPremium, ok = new(big.Int).SetString(jsonData.Fee.MaxCancellationPremium, 10)
		if !ok {
			return nil, errors.New("invalid maxCancellationPremium")
		}
	}
	resolverCancellationConfig, err := NewResolverCancellationConfig(maxCancellationPremium, jsonData.CancellationAuctionDuration)
	if err != nil {
		return nil, err
	}

	// Parse auction details
	points := make([]AuctionPoint, len(jsonData.DutchAuctionData.PointsAndTimeDeltas))
	for i, point := range jsonData.DutchAuctionData.PointsAndTimeDeltas {
		points[i] = AuctionPoint{
			Coefficient: point.RateBump,
			Delay:       point.TimeDelta,
		}
	}

	auctionDetails, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       jsonData.DutchAuctionData.StartTime,
		Duration:        jsonData.DutchAuctionData.Duration,
		InitialRateBump: jsonData.DutchAuctionData.InitialRateBump,
		Points:          points,
	})
	if err != nil {
		return nil, err
	}

	// Create order
	orderInfo := OrderInfoData{
		ID:                 jsonData.ID,
		SrcAmount:          srcAmount,
		MinDstAmount:       minDstAmount,
		EstimatedDstAmount: estimatedDstAmount,
		Receiver:           receiver,
		SrcMint:            srcMint,
		DstMint:            dstMint,
	}

	return NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *FeeConfig
		ResolverCancellationConfig *ResolverCancellationConfig
	}{
		SrcAssetIsNative:           jsonData.SrcAssetIsNative,
		DstAssetIsNative:           jsonData.DstAssetIsNative,
		OrderExpirationDelay:       orderExpirationDelay,
		Fees:                       fees,
		ResolverCancellationConfig: resolverCancellationConfig,
	})
}
