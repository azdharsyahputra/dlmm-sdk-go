package onchain

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// StaticParameters represents static parameters set by the protocol.
type StaticParameters struct {
	BaseFactor               uint16
	FilterPeriod             uint16
	DecayPeriod              uint16
	ReductionFactor          uint16
	VariableFeeControl       uint32
	MaxVolatilityAccumulator uint32
	MinBinId                 int32
	MaxBinId                 int32
	ProtocolShare            uint16
	BaseFeePowerFactor       uint8
	FunctionType             uint8
	CollectFeeMode           uint8
	Padding                  [3]uint8
}

// VariableParameters represents parameters that change based on market dynamics.
type VariableParameters struct {
	VolatilityAccumulator uint32
	VolatilityReference   uint32
	IndexReference        int32
	Padding               [4]uint8
	LastUpdateTimestamp   int64
	Padding1              [8]uint8
}

// ProtocolFee represents uncollected protocol fees.
type ProtocolFee struct {
	AmountX uint64
	AmountY uint64
}

// RewardInfo stores the state relevant for tracking liquidity mining rewards.
type RewardInfo struct {
	Mint                                     solana.PublicKey
	Vault                                    solana.PublicKey
	Funder                                   solana.PublicKey
	RewardDuration                           uint64
	RewardDurationEnd                        uint64
	RewardRate                               [16]byte // u128
	LastUpdateTime                           uint64
	CumulativeSecondsWithEmptyLiquidityReward uint64
}

// LbPair represents the LB pair state.
type LbPair struct {
	Parameters              StaticParameters
	VParameters             VariableParameters
	BumpSeed                [1]uint8
	BinStepSeed             [2]uint8
	PairType                uint8
	ActiveId                int32
	BinStep                 uint16
	Status                  uint8
	RequireBaseFactorSeed   uint8
	BaseFactorSeed          [2]uint8
	ActivationType          uint8
	CreatorPoolOnOffControl uint8
	TokenXMint              solana.PublicKey
	TokenYMint              solana.PublicKey
	ReserveX                solana.PublicKey
	ReserveY                solana.PublicKey
	ProtocolFee             ProtocolFee
	Padding1                [32]uint8
	RewardInfos             [2]RewardInfo
	Oracle                  solana.PublicKey
	BinArrayBitmap          [16]uint64
	LastUpdatedAt           int64
	Padding2                [32]uint8
	PreActivationSwapAddress solana.PublicKey
	BaseKey                 solana.PublicKey
	ActivationPoint         uint64
	PreActivationDuration   uint64
	Padding3                [8]uint8
	Padding4                uint64
	Creator                 solana.PublicKey
	TokenMintXProgramFlag   uint8
	TokenMintYProgramFlag   uint8
	Version                 uint8
	Reserved                [21]uint8
}

// Bin represents a single bin's state.
type Bin struct {
	AmountX                      uint64
	AmountY                      uint64
	Price                        [16]byte // u128
	LiquiditySupply              [16]byte // u128
	FulfilledOrderAmountX        uint64
	FulfilledOrderAmountY        uint64
	LimitOrderFeeAskSide         uint64
	LimitOrderFeeBidSide         uint64
	FeeAmountXPerTokenStored     [16]byte // u128
	FeeAmountYPerTokenStored     [16]byte // u128
	OpenOrderAmount              uint64
	TotalProcessingOrderAmount   uint64
	ProcessedOrderRemainingAmount uint64
	OrderAge                     uint32
	LimitOrderAskSide            uint8
	Padding1                     [3]uint8
}

// BinArray represents a range of bins.
type BinArray struct {
	Index    int64
	Version  uint8
	Padding1 [7]uint8
	LbPair   solana.PublicKey
	Bins     [70]Bin
}

// UserRewardInfo represents reward status for a user position in bins.
type UserRewardInfo struct {
	RewardPerTokenCompletes [2][16]byte // array of 2 u128
	RewardPendings          [2]uint64
}

// FeeInfo represents claimed fee status for a user position in bins.
type FeeInfo struct {
	FeeXPerTokenComplete [16]byte // u128
	FeeYPerTokenComplete [16]byte // u128
	FeeXPending          uint64
	FeeYPending          uint64
}

// PositionV2 represents a user position in DLMM.
type PositionV2 struct {
	LbPair                     solana.PublicKey
	Owner                      solana.PublicKey
	LiquidityShares            [70][16]byte // array of 70 u128
	RewardInfos                [70]UserRewardInfo
	FeeInfos                   [70]FeeInfo
	LowerBinId                 int32
	UpperBinId                 int32
	LastUpdatedAt              int64
	TotalClaimedFeeXAmount     uint64
	TotalClaimedFeeYAmount     uint64
	TotalClaimedRewards        [2]uint64
	Operator                   solana.PublicKey
	LockReleasePoint           uint64
	Padding0                   uint8
	FeeOwner                   solana.PublicKey
	Version                    uint8
	PermissionlessOperationBits uint8
	Reserved                   [85]uint8
}

// BinArrayBitmapExtension represents extension state when bitmap overflows.
type BinArrayBitmapExtension struct {
	LbPair                 solana.PublicKey
	PositiveBinArrayBitmap [12][8]uint64
	NegativeBinArrayBitmap [12][8]uint64
}

// DeserializeAccount deserializes Anchor account skipping the 8-byte discriminator.
func DeserializeAccount(data []byte, dest interface{}) error {
	if len(data) < 8 {
		return fmt.Errorf("account data too short (less than 8 bytes)")
	}
	// Skip 8-byte Anchor discriminator
	reader := bytes.NewReader(data[8:])
	if err := binary.Read(reader, binary.LittleEndian, dest); err != nil {
		return fmt.Errorf("failed to deserialize account: %w", err)
	}
	return nil
}
