package onchain

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
)

// FunctionType constants matching the TS SDK enum.
const (
	FunctionTypeUndetermined    uint8 = 0
	FunctionTypeLiquidityMining uint8 = 1
	FunctionTypeLimitOrder      uint8 = 2
)

// SwapQuote represents the simulated result of a swap.
type SwapQuote struct {
	InAmount              *big.Int
	OutAmount             *big.Int
	MinOutAmount          *big.Int // Minimum output amount after slippage
	Fee                   *big.Int
	ProtocolFee           *big.Int
	PriceImpact           float64
	LastFilledActiveBinId int32
	BinArraysPubkey       []solana.PublicKey // Bin array accounts traversed during the swap
}

// ComputeSwapQuote simulates a swap exact-in against a slice of BinArrays.
// Matches the TS SDK swapQuote() in index.ts:5527-5745.
func ComputeSwapQuote(
	lbPair *LbPair,
	binArrays []BinArray,
	inAmount *big.Int,
	swapForY bool,
	isPartialFill bool,
	allowedSlippageBps uint16,
	currentTimestamp int64,
	inTransferFee *TransferFee,
	outTransferFee *TransferFee,
) (*SwapQuote, error) {
	if inAmount == nil || inAmount.Sign() <= 0 {
		return nil, fmt.Errorf("input amount must be greater than zero")
	}

	if currentTimestamp == 0 {
		currentTimestamp = time.Now().Unix()
	}

	// Index bin arrays by their index
	binArrayMap := make(map[int64]BinArray)
	for _, ba := range binArrays {
		binArrayMap[ba.Index] = ba
	}

	transferFeeExcludedAmountIn, _ := CalculateTransferFeeExcludedAmount(inTransferFee, inAmount)
	inAmountLeft := new(big.Int).Set(transferFeeExcludedAmountIn)
	vParamClone := lbPair.VParameters
	activeId := lbPair.ActiveId

	binStep := lbPair.BinStep
	sParameters := lbParamToStaticParameters(lbPair)
	supportLimitOrder := isSupportLimitOrder(lbPair)
	feeOnInput := GetFeeMode(sParameters.CollectFeeMode, swapForY).FeeOnInput

	// 1. Decay/update reference before starting swap
	UpdateReference(activeId, &vParamClone, sParameters, currentTimestamp)

	var startBin *Bin
	var lastFilledActiveBinId int32 = activeId
	totalOutAmount := big.NewInt(0)
	feeAmount := big.NewInt(0)
	protocolFeeAmount := big.NewInt(0)
	binArraysForSwap := make(map[int64]bool)

	for inAmountLeft.Sign() > 0 {
		binArrayIndex := BinIdToBinArrayIndex(activeId)
		binArray, exists := binArrayMap[binArrayIndex]
		if !exists {
			if isPartialFill {
				break
			}
			return nil, fmt.Errorf("insufficient liquidity: BinArray at index %d not provided", binArrayIndex)
		}

		binArraysForSwap[binArrayIndex] = true

		lowerBinId, upperBinId := GetBinArrayLowerUpperBinId(binArrayIndex)
		if activeId < lowerBinId || activeId > upperBinId {
			return nil, fmt.Errorf("active ID %d out of bounds for BinArray index %d", activeId, binArrayIndex)
		}

		binIdx := activeId - lowerBinId
		bin := binArray.Bins[binIdx]

		maxAmountOut := GetBinMaxAmountOut(bin, swapForY, supportLimitOrder)
		if maxAmountOut.Sign() > 0 {
			UpdateVolatilityAccumulator(&vParamClone, sParameters, activeId)

			res := SwapExactInQuoteAtBin(
				bin,
				binStep,
				sParameters,
				vParamClone,
				inAmountLeft,
				swapForY,
				supportLimitOrder,
				feeOnInput,
			)

			if res.AmountIn.Sign() > 0 {
				inAmountLeft.Sub(inAmountLeft, res.AmountIn)
				totalOutAmount.Add(totalOutAmount, res.AmountOut)
				feeAmount.Add(feeAmount, res.Fee)
				protocolFeeAmount.Add(protocolFeeAmount, res.ProtocolFee)

				if startBin == nil {
					startBin = &bin
				}
				lastFilledActiveBinId = activeId
			}
		}

		if inAmountLeft.Sign() > 0 {
			if swapForY {
				activeId--
			} else {
				activeId++
			}
		}
	}

	if startBin == nil {
		return nil, fmt.Errorf("insufficient liquidity in the pool")
	}

	actualInAmount := new(big.Int).Sub(transferFeeExcludedAmountIn, inAmountLeft)

	transferFeeIncludedInAmount, _ := CalculateTransferFeeIncludedAmount(inTransferFee, actualInAmount)
	if transferFeeIncludedInAmount.Cmp(inAmount) > 0 {
		transferFeeIncludedInAmount = new(big.Int).Set(inAmount)
	}

	// Calculate price impact using computeFeeFromAmount (matches TS index.ts:5677-5694)
	feeForPriceImpact := ComputeFeeFromAmount(binStep, sParameters, vParamClone, actualInAmount)
	actualInAmountWithoutFees := new(big.Int).Sub(actualInAmount, feeForPriceImpact)
	startBinPrice := U128ToBigInt(startBin.Price)
	outAmountWithoutSlippage := GetAmountOut(startBinPrice, actualInAmountWithoutFees, swapForY, false)

	var priceImpact float64 = 0.0
	if outAmountWithoutSlippage.Sign() > 0 {
		outBig := new(big.Float).SetInt(totalOutAmount)
		outWithoutSlippageBig := new(big.Float).SetInt(outAmountWithoutSlippage)

		diff := new(big.Float).Sub(outBig, outWithoutSlippageBig)
		ratio := new(big.Float).Quo(diff, outWithoutSlippageBig)
		ratio.Mul(ratio, big.NewFloat(100.0))

		impactVal, _ := ratio.Float64()
		priceImpact = math.Abs(impactVal)
	}

	transferFeeExcludedAmountOut, _ := CalculateTransferFeeExcludedAmount(outTransferFee, totalOutAmount)

	// Calculate minOutAmount after allowed slippage
	factor := big.NewInt(int64(10000 - allowedSlippageBps))
	minOutAmount := new(big.Int).Mul(transferFeeExcludedAmountOut, factor)
	minOutAmount.Div(minOutAmount, big.NewInt(10000))

	// Collect bin array pubkeys that were traversed
	var binArraysPubkey []solana.PublicKey
	for _, ba := range binArrays {
		if binArraysForSwap[ba.Index] {
			pubkey, _, err := DeriveBinArray(lbPair.ReserveX, ba.Index) // derive from pool address
			if err == nil {
				binArraysPubkey = append(binArraysPubkey, pubkey)
			}
		}
	}

	return &SwapQuote{
		InAmount:              transferFeeIncludedInAmount,
		OutAmount:             transferFeeExcludedAmountOut,
		MinOutAmount:          minOutAmount,
		Fee:                   feeAmount,
		ProtocolFee:           protocolFeeAmount,
		PriceImpact:           priceImpact,
		LastFilledActiveBinId: lastFilledActiveBinId,
		BinArraysPubkey:       binArraysPubkey,
	}, nil
}

// lbParamToStaticParameters wraps LbPair parameters to StaticParameters.
func lbParamToStaticParameters(lbPair *LbPair) StaticParameters {
	return lbPair.Parameters
}

// isSupportLimitOrder checks if the pool supports limit orders.
// Matches TS SDK lbPair.ts:46-60 exactly.
func isSupportLimitOrder(lbPair *LbPair) bool {
	functionType := lbPair.Parameters.FunctionType
	switch functionType {
	case FunctionTypeLimitOrder:
		return true
	case FunctionTypeLiquidityMining:
		return false
	case FunctionTypeUndetermined:
		// Support limit orders only if no rewards are initialized
		// (reward mint is zero/default pubkey)
		for _, ri := range lbPair.RewardInfos {
			if !ri.Mint.IsZero() {
				return false
			}
		}
		return true
	default:
		return false
	}
}
