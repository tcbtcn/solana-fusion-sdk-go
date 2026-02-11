package quoter

// PresetDTO represents a preset configuration
type PresetDTO struct {
	StartAuctionIn     uint32           `json:"startAuctionIn"`
	AuctionDuration    uint32           `json:"auctionDuration"`
	InitialRateBump    uint16           `json:"initialRateBump"`
	AuctionStartAmount string           `json:"auctionStartAmount"`
	AuctionEndAmount   string           `json:"auctionEndAmount"`
	CostInDstToken     string           `json:"costInDstToken"`
	Points             []PresetPointDTO `json:"points"`
}

// PresetPointDTO represents a point in a preset
type PresetPointDTO struct {
	Delay       uint16 `json:"delay"`
	Coefficient uint16 `json:"coefficient"`
}

// PresetType represents the type of preset
type PresetType string

const (
	PresetTypeFast   PresetType = "fast"
	PresetTypeMedium PresetType = "medium"
	PresetTypeSlow   PresetType = "slow"
)

// QuoteDTO represents a quote response
type QuoteDTO struct {
	QuoteID            *string    `json:"quoteId"`
	SrcAmount          string     `json:"srcAmount"`
	DstAmount          string     `json:"dstAmount"`
	Presets            PresetsDTO `json:"presets"`
	RecommendedPreset  PresetType `json:"recommendedPreset"`
	Prices             PricesDTO  `json:"prices"`
	Volume             VolumeDTO  `json:"volume"`
	PriceImpactPercent float64    `json:"priceImpactPercent"`
}

// PresetsDTO represents all presets
type PresetsDTO struct {
	Fast   PresetDTO `json:"fast"`
	Medium PresetDTO `json:"medium"`
	Slow   PresetDTO `json:"slow"`
}

// PricesDTO represents price information
type PricesDTO struct {
	USD USDPriceDTO `json:"usd"`
}

// USDPriceDTO represents USD price information
type USDPriceDTO struct {
	SrcToken string `json:"srcToken"`
	DstToken string `json:"dstToken"`
}

// VolumeDTO represents volume information
type VolumeDTO struct {
	USD USDVolumeDTO `json:"usd"`
}

// USDVolumeDTO represents USD volume information
type USDVolumeDTO struct {
	SrcToken string `json:"srcToken"`
	DstToken string `json:"dstToken"`
}
