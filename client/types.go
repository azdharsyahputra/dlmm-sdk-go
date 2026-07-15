package client

// TokenMetrics represents metrics for a token in the liquidity pair.
type TokenMetrics struct {
	Address                 string  `json:"address"`
	Decimals                int32   `json:"decimals"`
	FreezeAuthorityDisabled bool    `json:"freeze_authority_disabled"`
	Holders                 int32   `json:"holders"`
	IsVerified              bool    `json:"is_verified"`
	MarketCap               float64 `json:"market_cap"`
	Name                    string  `json:"name"`
	Price                   float64 `json:"price"`
	Symbol                  string  `json:"symbol"`
	TotalSupply             float64 `json:"total_supply"`
}

// PoolConfig represents config of the pool.
type PoolConfig struct {
	BaseFeePct     float64 `json:"base_fee_pct"`
	BinStep        int32   `json:"bin_step"`
	CollectFeeMode int32   `json:"collect_fee_mode"`
	MaxFeePct      float64 `json:"max_fee_pct"`
	ProtocolFeePct float64 `json:"protocol_fee_pct"`
}

// TimeWindowData represents volume/fee data over different windows.
type TimeWindowData struct {
	M30 float64 `json:"30m"`
	H1  float64 `json:"1h"`
	H2  float64 `json:"2h"`
	H4  float64 `json:"4h"`
	H12 float64 `json:"12h"`
	H24 float64 `json:"24h"`
}

// CumulativeMetrics represents cumulative volume and fees.
type CumulativeMetrics struct {
	Fees   float64 `json:"fees"`
	Volume float64 `json:"volume"`
}

// PoolInfo represents metadata and state of a single pool.
type PoolInfo struct {
	Address           string            `json:"address"`
	Name              string            `json:"name"`
	TokenX            TokenMetrics      `json:"token_x"`
	TokenY            TokenMetrics      `json:"token_y"`
	ReserveX          string            `json:"reserve_x"`
	ReserveY          string            `json:"reserve_y"`
	TokenXAmount      float64           `json:"token_x_amount"`
	TokenYAmount      float64           `json:"token_y_amount"`
	CreatedAt         int64             `json:"created_at"`
	RewardMintX       string            `json:"reward_mint_x"`
	RewardMintY       string            `json:"reward_mint_y"`
	PoolConfig        PoolConfig        `json:"pool_config"`
	DynamicFeePct     float64           `json:"dynamic_fee_pct"`
	Tvl               float64           `json:"tvl"`
	CurrentPrice      float64           `json:"current_price"`
	Apr               float64           `json:"apr"`
	Apy               float64           `json:"apy"`
	HasFarm           bool              `json:"has_farm"`
	FarmApr           float64           `json:"farm_apr"`
	FarmApy           float64           `json:"farm_apy"`
	Volume            TimeWindowData    `json:"volume"`
	Fees              TimeWindowData    `json:"fees"`
	ProtocolFees      TimeWindowData    `json:"protocol_fees"`
	FeeTvlRatio       TimeWindowData    `json:"fee_tvl_ratio"`
	CumulativeMetrics CumulativeMetrics `json:"cumulative_metrics"`
	IsBlacklisted     bool              `json:"is_blacklisted"`
	Tags              []string          `json:"tags"`
	Launchpad         *string           `json:"launchpad,omitempty"`
}

// PoolsResponse represents paginated response of pools.
type PoolsResponse struct {
	CurrentPage int64      `json:"current_page"`
	PageSize    int64      `json:"page_size"`
	Pages       int64      `json:"pages"`
	Total       int64      `json:"total"`
	Data        []PoolInfo `json:"data"`
}

// PoolResponse is returned by /pools/{address}.
type PoolResponse PoolInfo

// OHLCVItem represents a candle for a single pool.
type OHLCVItem struct {
	Timestamp    int64   `json:"timestamp"`
	TimestampStr string  `json:"timestamp_str"`
	Open         float64 `json:"open"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Close        float64 `json:"close"`
	Volume       float64 `json:"volume"`
}

// OHLCVResponse represents OHLCV candles response.
type OHLCVResponse struct {
	StartTime int64       `json:"start_time"`
	EndTime   int64       `json:"end_time"`
	Data      []OHLCVItem `json:"data"`
}

// VolumeHistoryItem represents historical volume for a pool.
type VolumeHistoryItem struct {
	Timestamp    int64   `json:"timestamp"`
	TimestampStr string  `json:"timestamp_str"`
	Volume       float64 `json:"volume"`
	Fees         float64 `json:"fees"`
	ProtocolFees float64 `json:"protocol_fees"`
}

// VolumeHistoryResponse represents volume history response.
type VolumeHistoryResponse struct {
	StartTime int64               `json:"start_time"`
	EndTime   int64               `json:"end_time"`
	Data      []VolumeHistoryItem `json:"data"`
}

// TokenAmount represents token quantity in amount, value in SOL, and value in USD.
type TokenAmount struct {
	Amount    string  `json:"amount"`
	AmountSol *string `json:"amountSol,omitempty"`
	AmountUsd string  `json:"amountUsd"`
}

// TokenPairWithTotal represents per-token amounts with a combined total.
type TokenPairWithTotal struct {
	TokenX TokenAmount `json:"tokenX"`
	TokenY TokenAmount `json:"tokenY"`
	Total  TokenAmount `json:"total"`
}

// PoolPortfolioItem represents portfolio item for a pool.
type PoolPortfolioItem struct {
	PoolAddress              string `json:"poolAddress"`
	BinStep                  string `json:"binStep"`
	BaseFee                  string `json:"baseFee"`
	CollectFeeMode           int32  `json:"collectFeeMode"`
	TokenXMint               string `json:"tokenXMint"`
	TokenYMint               string `json:"tokenYMint"`
	TokenXIcon               string `json:"tokenXIcon"`
	TokenYIcon               string `json:"tokenYIcon"`
	TokenX                   string `json:"tokenX"`
	TokenY                   string `json:"tokenY"`
	TotalDeposit             string `json:"totalDeposit"`
	TotalDepositSol          string `json:"totalDepositSol"`
	TotalWithdrawal          string `json:"totalWithdrawal"`
	TotalWithdrawalSol       string `json:"totalWithdrawalSol"`
	TotalFee                 string `json:"totalFee"`
	TotalFeeSol              string `json:"totalFeeSol"`
	PnlUsd                   string `json:"pnlUsd"`
	PnlSol                   string `json:"pnlSol"`
	PnlPctChange             string `json:"pnlPctChange"`
	PnlSolPctChange          string `json:"pnlSolPctChange"`
	TotalDepositTokenX       string `json:"totalDepositTokenX"`
	TotalDepositTokenXUsd    string `json:"totalDepositTokenXUsd"`
	TotalDepositTokenXSol    string `json:"totalDepositTokenXSol"`
	TotalWithdrawalTokenX    string `json:"totalWithdrawalTokenX"`
	TotalWithdrawalTokenXUsd string `json:"totalWithdrawalTokenXUsd"`
	TotalWithdrawalTokenXSol string `json:"totalWithdrawalTokenXSol"`
	TotalFeeTokenX           string `json:"totalFeeTokenX"`
	TotalFeeTokenXUsd        string `json:"totalFeeTokenXUsd"`
	TotalFeeTokenXSol        string `json:"totalFeeTokenXSol"`
	TotalDepositTokenY       string `json:"totalDepositTokenY"`
	TotalDepositTokenYUsd    string `json:"totalDepositTokenYUsd"`
	TotalDepositTokenYSol    string `json:"totalDepositTokenYSol"`
	TotalWithdrawalTokenY    string `json:"totalWithdrawalTokenY"`
	TotalWithdrawalTokenYUsd string `json:"totalWithdrawalTokenYUsd"`
	TotalWithdrawalTokenYSol string `json:"totalWithdrawalTokenYSol"`
	TotalFeeTokenY           string `json:"totalFeeTokenY"`
	TotalFeeTokenYUsd        string `json:"totalFeeTokenYUsd"`
	TotalFeeTokenYSol        string `json:"totalFeeTokenYSol"`
}

// PortfolioResponse represents response for closed portfolio.
type PortfolioResponse struct {
	HasNext        bool                `json:"hasNext"`
	Page           int32               `json:"page"`
	PageSize       int32               `json:"pageSize"`
	Pools          []PoolPortfolioItem `json:"pools"`
	TotalCount     int64               `json:"totalCount"`
	TotalPositions int64               `json:"totalPositions"`
}

// PoolOpenPortfolioItem represents open portfolio details for a single pool.
type PoolOpenPortfolioItem struct {
	Balances            string   `json:"balances"`
	BalancesSol         *string  `json:"balancesSol,omitempty"`
	BaseFee             float64  `json:"baseFee"`
	BinStep             int32    `json:"binStep"`
	CollectFeeMode      int32    `json:"collectFeeMode"`
	FeePerTvl24h        string   `json:"feePerTvl24h"`
	ListPositions       []string `json:"listPositions"`
	OpenPositionCount   int64    `json:"openPositionCount"`
	OutOfRange          *bool    `json:"outOfRange,omitempty"`
	Pnl                 string   `json:"pnl"`
	PnlPctChange        string   `json:"pnlPctChange"`
	PnlSol              *string  `json:"pnlSol,omitempty"`
	PnlSolPctChange     *string  `json:"pnlSolPctChange,omitempty"`
	PoolAddress         string   `json:"poolAddress"`
	PositionsOutOfRange int64    `json:"positionsOutOfRange"`
	RewardX             *string  `json:"rewardX,omitempty"`
	RewardY             *string  `json:"rewardY,omitempty"`
	TokenX              string   `json:"tokenX"`
	TokenXIcon          string   `json:"tokenXIcon"`
	TokenXMint          string   `json:"tokenXMint"`
	TokenY              string   `json:"tokenY"`
	TokenYIcon          string   `json:"tokenYIcon"`
	TokenYMint          string   `json:"tokenYMint"`
	TotalDeposit        string   `json:"totalDeposit"`
	TotalDepositSol     *string  `json:"totalDepositSol,omitempty"`
	UnclaimedFees       string   `json:"unclaimedFees"`
	UnclaimedFeesSol    *string  `json:"unclaimedFeesSol,omitempty"`
}

// OpenPortfolioResponse represents open portfolio details.
type OpenPortfolioResponse struct {
	HasNext        bool                    `json:"hasNext"`
	Page           int32                   `json:"page"`
	PageSize       int32                   `json:"pageSize"`
	Pools          []PoolOpenPortfolioItem `json:"pools"`
	TotalCount     int64                   `json:"totalCount"`
	TotalPositions int64                   `json:"totalPositions"`
	SolPrice       *string                 `json:"solPrice,omitempty"`
}

// PortfolioTotalResponse represents total portfolio metadata.
type PortfolioTotalResponse struct {
	TotalClosedPositions int64  `json:"totalClosedPositions"`
	TotalPnlPctChange    string `json:"totalPnlPctChange"`
	TotalPnlSol          string `json:"totalPnlSol"`
	TotalPnlSolPctChange string `json:"totalPnlSolPctChange"`
	TotalPnlUsd          string `json:"totalPnlUsd"`
}

// PositionPnLData represents PnL details for a single position.
type PositionPnLData struct {
	PositionAddress    string             `json:"positionAddress"`
	MinPrice           string             `json:"minPrice"`
	MaxPrice           string             `json:"maxPrice"`
	LowerBinId         int32              `json:"lowerBinId"`
	UpperBinId         int32              `json:"upperBinId"`
	FeePerTvl24h       string             `json:"feePerTvl24h"`
	IsClosed           bool               `json:"isClosed"`
	IsOutOfRange       *bool              `json:"isOutOfRange,omitempty"`
	PnlUsd             string             `json:"pnlUsd"`
	PnlPctChange       string             `json:"pnlPctChange"`
	PnlSol             *string            `json:"pnlSol,omitempty"`
	PnlSolPctChange    *string            `json:"pnlSolPctChange,omitempty"`
	AllTimeDeposits    TokenPairWithTotal `json:"allTimeDeposits"`
	AllTimeWithdrawals TokenPairWithTotal `json:"allTimeWithdrawals"`
	AllTimeFees        TokenPairWithTotal `json:"allTimeFees"`
	CreatedAt          *int64             `json:"createdAt,omitempty"`
	ClosedAt           *int64             `json:"closedAt,omitempty"`
}

// GetPoolPositionPnLResponse represents response from pool position PnL.
type GetPoolPositionPnLResponse struct {
	TotalCount        int64             `json:"totalCount"`
	Page              int32             `json:"page"`
	PageSize          int32             `json:"pageSize"`
	HasNext           bool              `json:"hasNext"`
	Positions         []PositionPnLData `json:"positions"`
	TokenX            *string           `json:"tokenX,omitempty"`
	TokenY            *string           `json:"tokenY,omitempty"`
	RewardTokenX      *string           `json:"rewardTokenX,omitempty"`
	RewardTokenY      *string           `json:"rewardTokenY,omitempty"`
	TokenXPrice       string            `json:"tokenXPrice"`
	TokenYPrice       string            `json:"tokenYPrice"`
	RewardTokenXPrice string            `json:"rewardTokenXPrice"`
	RewardTokenYPrice string            `json:"rewardTokenYPrice"`
	SolPrice          *string           `json:"solPrice,omitempty"`
}

// PositionEvent represents a single position action event.
type PositionEvent struct {
	EventAddress string      `json:"eventAddress"`
	EventType    string      `json:"eventType"`
	TxSignature  string      `json:"txSignature"`
	Slot         int64       `json:"slot"`
	BlockTime    int64       `json:"blockTime"`
	TokenXAmount string      `json:"tokenXAmount"`
	TokenYAmount string      `json:"tokenYAmount"`
	TokenXUsd    string      `json:"tokenXUsd"`
	TokenYUsd    string      `json:"tokenYUsd"`
	TotalUsd     string      `json:"totalUsd"`
	TokenXSol    *string     `json:"tokenXSol,omitempty"`
	TokenYSol    *string     `json:"tokenYSol,omitempty"`
	TotalSol     *string     `json:"totalSol,omitempty"`
	Rewards      TokenAmount `json:"rewards"`
}

// GetPositionHistoricalEventsResponse represents historical events.
type GetPositionHistoricalEventsResponse struct {
	Data []PositionEvent `json:"data"`
}

// DailyMetricPoint represents a daily metrics data point.
type DailyMetricPoint struct {
	Timestamp int64    `json:"timestamp"`
	DateTime  string   `json:"date_time"`
	Clean     *float64 `json:"clean,omitempty"`
	Gross     *float64 `json:"gross,omitempty"`
}

// DailyMetricResponse represents daily metrics response.
type DailyMetricResponse struct {
	DataPoints  []DailyMetricPoint `json:"data_points"`
	RefreshedAt *int64             `json:"refreshed_at,omitempty"`
}

// ProtocolMetricsResponse represents protocol global metrics.
type ProtocolMetricsResponse struct {
	Tvl            float64 `json:"tvl"`
	Volume24h      float64 `json:"volume_24h"`
	Fees24h        float64 `json:"fees_24h"`
	TradingFees24h float64 `json:"trading_fees_24h"`
}

// LimitOrderPoolDetails represents pool config for limit orders.
type LimitOrderPoolDetails struct {
	Address        string `json:"address"`
	Name           string `json:"name"`
	TokenX         string `json:"token_x"`
	TokenY         string `json:"token_y"`
	TokenXMint     string `json:"token_x_mint"`
	TokenYMint     string `json:"token_y_mint"`
	TokenXDecimals int32  `json:"token_x_decimals"`
	TokenYDecimals int32  `json:"token_y_decimals"`
	BaseFeePct     string `json:"base_fee_pct"`
	BinStep        int32  `json:"bin_step"`
}

// OpenLimitOrderPoolSummary represents summary of open orders in a pool.
type OpenLimitOrderPoolSummary struct {
	Pool                LimitOrderPoolDetails `json:"pool"`
	TotalOrders         int64                 `json:"total_orders"`
	FilledPct           string                `json:"filled_pct"`
	MinLimitPrice       string                `json:"min_limit_price"`
	MaxLimitPrice       string                `json:"max_limit_price"`
	TotalDepositX       string                `json:"total_deposit_x"`
	TotalDepositY       string                `json:"total_deposit_y"`
	TotalDepositXUsd    string                `json:"total_deposit_x_usd"`
	TotalDepositXSol    string                `json:"total_deposit_x_sol"`
	TotalDepositYUsd    string                `json:"total_deposit_y_usd"`
	TotalDepositYSol    string                `json:"total_deposit_y_sol"`
	TotalDepositUsd     string                `json:"total_deposit_usd"`
	TotalDepositSol     string                `json:"total_deposit_sol"`
	UnrealizedPnlUsd    string                `json:"unrealized_pnl_usd"`
	UnrealizedPnlSol    string                `json:"unrealized_pnl_sol"`
	UnrealizedPnlPctUsd string                `json:"unrealized_pnl_pct_usd"`
	UnrealizedPnlPctSol string                `json:"unrealized_pnl_pct_sol"`
}

// PaginationResponse_OpenLimitOrderPoolSummary represents open limit order pools response.
type PaginationResponse_OpenLimitOrderPoolSummary struct {
	CurrentPage int64                       `json:"current_page"`
	PageSize    int64                       `json:"page_size"`
	Pages       int64                       `json:"pages"`
	Total       int64                       `json:"total"`
	Data        []OpenLimitOrderPoolSummary `json:"data"`
}

// ClosedLimitOrderPoolSummary represents summary of closed orders in a pool.
type ClosedLimitOrderPoolSummary struct {
	Pool                LimitOrderPoolDetails `json:"pool"`
	TotalOrders         int64                 `json:"total_orders"`
	FullyFilledOrders   int64                 `json:"fully_filled_orders"`
	FilledPct           string                `json:"filled_pct"`
	MinLimitPrice       string                `json:"min_limit_price"`
	MaxLimitPrice       string                `json:"max_limit_price"`
	TotalDepositX       string                `json:"total_deposit_x"`
	TotalDepositY       string                `json:"total_deposit_y"`
	TotalDepositXUsd    string                `json:"total_deposit_x_usd"`
	TotalDepositXSol    string                `json:"total_deposit_x_sol"`
	TotalDepositYUsd    string                `json:"total_deposit_y_usd"`
	TotalDepositYSol    string                `json:"total_deposit_y_sol"`
	TotalDepositUsd     string                `json:"total_deposit_usd"`
	TotalDepositSol     string                `json:"total_deposit_sol"`
	TotalWithdrawalX    string                `json:"total_withdrawal_x"`
	TotalWithdrawalY    string                `json:"total_withdrawal_y"`
	TotalWithdrawalXUsd string                `json:"total_withdrawal_x_usd"`
	TotalWithdrawalXSol string                `json:"total_withdrawal_x_sol"`
	TotalWithdrawalYUsd string                `json:"total_withdrawal_y_usd"`
	TotalWithdrawalYSol string                `json:"total_withdrawal_y_sol"`
	TotalWithdrawalUsd  string                `json:"total_withdrawal_usd"`
	TotalWithdrawalSol  string                `json:"total_withdrawal_sol"`
	TotalBonusX         string                `json:"total_bonus_x"`
	TotalBonusY         string                `json:"total_bonus_y"`
	TotalBonusXUsd      string                `json:"total_bonus_x_usd"`
	TotalBonusXSol      string                `json:"total_bonus_x_sol"`
	TotalBonusYUsd      string                `json:"total_bonus_y_usd"`
	TotalBonusYSol      string                `json:"total_bonus_y_sol"`
	TotalBonusUsd       string                `json:"total_bonus_usd"`
	TotalBonusSol       string                `json:"total_bonus_sol"`
	RealizedPnlUsd      string                `json:"realized_pnl_usd"`
	RealizedPnlSol      string                `json:"realized_pnl_sol"`
	RealizedPnlPctUsd   string                `json:"realized_pnl_pct_usd"`
	RealizedPnlPctSol   string                `json:"realized_pnl_pct_sol"`
	LastClosedAt        int64                 `json:"last_closed_at"`
}

// PaginationResponse_ClosedLimitOrderPoolSummary represents closed limit order pools response.
type PaginationResponse_ClosedLimitOrderPoolSummary struct {
	CurrentPage int64                         `json:"current_page"`
	PageSize    int64                         `json:"page_size"`
	Pages       int64                         `json:"pages"`
	Total       int64                         `json:"total"`
	Data        []ClosedLimitOrderPoolSummary `json:"data"`
}

// LimitOrderSummaryResponse represents limit order totals.
type LimitOrderSummaryResponse struct {
	OpenOrdersCount     int64  `json:"open_orders_count"`
	ClosedOrdersCount   int64  `json:"closed_orders_count"`
	TotalDepositUsd     string `json:"total_deposit_usd"`
	TotalDepositSol     string `json:"total_deposit_sol"`
	TotalBonusUsd       string `json:"total_bonus_usd"`
	TotalBonusSol       string `json:"total_bonus_sol"`
	CurrentSpotValueUsd string `json:"current_spot_value_usd"`
	CurrentSpotValueSol string `json:"current_spot_value_sol"`
	RealizedPnlUsd      string `json:"realized_pnl_usd"`
	RealizedPnlSol      string `json:"realized_pnl_sol"`
	RealizedPnlPctUsd   string `json:"realized_pnl_pct_usd"`
	RealizedPnlPctSol   string `json:"realized_pnl_pct_sol"`
	UnrealizedPnlUsd    string `json:"unrealized_pnl_usd"`
	UnrealizedPnlSol    string `json:"unrealized_pnl_sol"`
	UnrealizedPnlPctUsd string `json:"unrealized_pnl_pct_usd"`
	UnrealizedPnlPctSol string `json:"unrealized_pnl_pct_sol"`
}

// GetWalletTotalClaimsResponse represents total claims for a wallet.
type GetWalletTotalClaimsResponse struct {
	ClaimedRewardX    string  `json:"claimedRewardX"`
	ClaimedRewardXUsd string  `json:"claimedRewardXUsd"`
	ClaimedRewardXSol *string `json:"claimedRewardXSol,omitempty"`
	ClaimedRewardY    string  `json:"claimedRewardY"`
	ClaimedRewardYUsd string  `json:"claimedRewardYUsd"`
	ClaimedRewardYSol *string `json:"claimedRewardYSol,omitempty"`
	ClaimedFeeX       string  `json:"claimedFeeX"`
	ClaimedFeeXUsd    string  `json:"claimedFeeXUsd"`
	ClaimedFeeXSol    *string `json:"claimedFeeXSol,omitempty"`
	ClaimedFeeY       string  `json:"claimedFeeY"`
	ClaimedFeeYUsd    string  `json:"claimedFeeYUsd"`
	ClaimedFeeYSol    *string `json:"claimedFeeYSol,omitempty"`
	TotalUsd          string  `json:"totalUsd"`
	TotalSol          *string `json:"totalSol,omitempty"`
}

// BonusClaimedResponse represents total realized bonus claimed on one pool.
type BonusClaimedResponse struct {
	CancelEventCount int64  `json:"cancel_event_count"`
	TotalBonusSol    string `json:"total_bonus_sol"`
	TotalBonusUsd    string `json:"total_bonus_usd"`
	TotalBonusX      string `json:"total_bonus_x"`
	TotalBonusXSol   string `json:"total_bonus_x_sol"`
	TotalBonusXUsd   string `json:"total_bonus_x_usd"`
	TotalBonusY      string `json:"total_bonus_y"`
	TotalBonusYSol   string `json:"total_bonus_y_sol"`
	TotalBonusYUsd   string `json:"total_bonus_y_usd"`
}

// ClosedLimitOrderDetail represents per-order lifecycle details for a closed order.
type ClosedLimitOrderDetail struct {
	LimitOrderAddress    string  `json:"limit_order_address"`
	UserAddress          string  `json:"user_address"`
	IsAskSide            bool    `json:"is_ask_side"`
	LowerPoolPrice       float64 `json:"lower_pool_price"`
	UpperPoolPrice       float64 `json:"upper_pool_price"`
	InputToken           string  `json:"input_token"`
	InputTokenMint       string  `json:"input_token_mint"`
	OutputToken          string  `json:"output_token"`
	OutputTokenMint      string  `json:"output_token_mint"`
	TotalDepositX        string  `json:"total_deposit_x"`
	TotalDepositY        string  `json:"total_deposit_y"`
	TotalDepositUsd      string  `json:"total_deposit_usd"`
	TotalDepositSol      string  `json:"total_deposit_sol"`
	OutputAmountExpected string  `json:"output_amount_expected"`
	TotalWithdrawalX     string  `json:"total_withdrawal_x"`
	TotalWithdrawalY     string  `json:"total_withdrawal_y"`
	TotalWithdrawalXUsd  string  `json:"total_withdrawal_x_usd"`
	TotalWithdrawalXSol  string  `json:"total_withdrawal_x_sol"`
	TotalWithdrawalYUsd  string  `json:"total_withdrawal_y_usd"`
	TotalWithdrawalYSol  string  `json:"total_withdrawal_y_sol"`
	TotalWithdrawalUsd   string  `json:"total_withdrawal_usd"`
	TotalWithdrawalSol   string  `json:"total_withdrawal_sol"`
	TotalBonusX          string  `json:"total_bonus_x"`
	TotalBonusY          string  `json:"total_bonus_y"`
	TotalBonusXUsd       string  `json:"total_bonus_x_usd"`
	TotalBonusXSol       string  `json:"total_bonus_x_sol"`
	TotalBonusYUsd       string  `json:"total_bonus_y_usd"`
	TotalBonusYSol       string  `json:"total_bonus_y_sol"`
	TotalBonusUsd        string  `json:"total_bonus_usd"`
	TotalBonusSol        string  `json:"total_bonus_sol"`
	RealizedPnlUsd       string  `json:"realized_pnl_usd"`
	RealizedPnlSol       string  `json:"realized_pnl_sol"`
	RealizedPnlPctUsd    string  `json:"realized_pnl_pct_usd"`
	RealizedPnlPctSol    string  `json:"realized_pnl_pct_sol"`
	FilledPct            string  `json:"filled_pct"`
	FilledInputAmount    string  `json:"filled_input_amount"`
	ReceivedOutputAmount string  `json:"received_output_amount"`
	OpenedAt             int64   `json:"opened_at"`
	OpenedAtSignature    string  `json:"opened_at_signature"`
	OpenedAtSlot         int64   `json:"opened_at_slot"`
	LastClosedAt         int64   `json:"last_closed_at"`
	TerminalSignature    string  `json:"terminal_signature"`
	TerminalSlot         int64   `json:"terminal_slot"`
}

// ClosedLimitOrdersForPoolResponse represents response for closed limit orders in a pool.
type ClosedLimitOrdersForPoolResponse struct {
	Pool        LimitOrderPoolDetails    `json:"pool"`
	Total       int64                    `json:"total"`
	Pages       int64                    `json:"pages"`
	CurrentPage int64                    `json:"current_page"`
	PageSize    int64                    `json:"page_size"`
	Data        []ClosedLimitOrderDetail `json:"data"`
}

// OpenLimitOrderBinPoint represents a bin description in an open limit order.
type OpenLimitOrderBinPoint struct {
	BinId                int32  `json:"bin_id"`
	Price                string `json:"price"`
	DepositAmount        string `json:"deposit_amount"`
	FulfilledAmount      string `json:"fulfilled_amount"`
	UnfilledAmount       string `json:"unfilled_amount"`
	OutputReceivedAmount string `json:"output_received_amount"`
	FillStatus           string `json:"fill_status"`
}

// OpenLimitOrderDetail represents per-order details for an open order.
type OpenLimitOrderDetail struct {
	LimitOrderAddress      string                   `json:"limit_order_address"`
	UserAddress            string                   `json:"user_address"`
	IsAskSide              bool                     `json:"is_ask_side"`
	LowerPoolPrice         float64                  `json:"lower_pool_price"`
	UpperPoolPrice         float64                  `json:"upper_pool_price"`
	InputToken             string                   `json:"input_token"`
	InputTokenMint         string                   `json:"input_token_mint"`
	OutputToken            string                   `json:"output_token"`
	OutputTokenMint        string                   `json:"output_token_mint"`
	InputAmount            string                   `json:"input_amount"`
	InputAmountUsd         string                   `json:"input_amount_usd"`
	InputAmountSol         string                   `json:"input_amount_sol"`
	OutputAmountExpected   string                   `json:"output_amount_expected"`
	FilledPct              string                   `json:"filled_pct"`
	FilledInputAmount      string                   `json:"filled_input_amount"`
	TotalFilledAmount      string                   `json:"total_filled_amount"`
	TotalFilledAmountUsd   string                   `json:"total_filled_amount_usd"`
	TotalFilledAmountSol   string                   `json:"total_filled_amount_sol"`
	TotalUnfilledAmount    string                   `json:"total_unfilled_amount"`
	TotalUnfilledAmountUsd string                   `json:"total_unfilled_amount_usd"`
	TotalUnfilledAmountSol string                   `json:"total_unfilled_amount_sol"`
	TotalBonusX            string                   `json:"total_bonus_x"`
	TotalBonusY            string                   `json:"total_bonus_y"`
	TotalBonusUsd          string                   `json:"total_bonus_usd"`
	TotalBonusSol          string                   `json:"total_bonus_sol"`
	UnrealizedPnlUsd       string                   `json:"unrealized_pnl_usd"`
	UnrealizedPnlSol       string                   `json:"unrealized_pnl_sol"`
	UnrealizedPnlPctUsd    string                   `json:"unrealized_pnl_pct_usd"`
	UnrealizedPnlPctSol    string                   `json:"unrealized_pnl_pct_sol"`
	OpenedAt               int64                    `json:"opened_at"`
	OpenedAtSignature      string                   `json:"opened_at_signature"`
	OpenedAtSlot           int64                    `json:"opened_at_slot"`
	BinDistribution        []OpenLimitOrderBinPoint `json:"bin_distribution"`
}

// OpenLimitOrdersForPoolResponse represents response for open limit orders in a pool.
type OpenLimitOrdersForPoolResponse struct {
	Pool               LimitOrderPoolDetails  `json:"pool"`
	CurrentActiveBinId int32                  `json:"current_active_bin_id"`
	CurrentPoolPrice   string                 `json:"current_pool_price"`
	Total              int64                  `json:"total"`
	Pages              int64                  `json:"pages"`
	CurrentPage        int64                  `json:"current_page"`
	PageSize           int64                  `json:"page_size"`
	Data               []OpenLimitOrderDetail `json:"data"`
}
