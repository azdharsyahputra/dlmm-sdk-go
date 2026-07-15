package onchain

import (
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

const BasisPointMax = 10000.0

// Constants for fixed-point math
var (
	One                = new(big.Int).Lsh(big.NewInt(1), 64) // 2^64
	FeePrecision       = big.NewInt(10_000_000_000)          // 10^10
	MaxFeeRate         = big.NewInt(1_000_000_000)           // 10%
	LimitOrderFeeShare = big.NewInt(5000)                    // 50%
	BasisPointMaxBig   = big.NewInt(10000)                   // 10000
)

// U128ToBigInt converts a little-endian [16]byte u128 to *big.Int.
func U128ToBigInt(u128 [16]byte) *big.Int {
	be := make([]byte, 16)
	for i := 0; i < 16; i++ {
		be[i] = u128[15-i]
	}
	return new(big.Int).SetBytes(be)
}

// BigIntToU128 converts *big.Int to a little-endian [16]byte u128.
func BigIntToU128(val *big.Int) [16]byte {
	var out [16]byte
	b := val.Bytes()
	// b is big-endian
	for i := 0; i < len(b) && i < 16; i++ {
		out[i] = b[len(b)-1-i]
	}
	return out
}

// MulShr performs: (x * y) >> offset, optionally rounding up if there is a remainder.
func MulShr(x, y *big.Int, offset uint, roundUp bool) *big.Int {
	prod := new(big.Int).Mul(x, y)
	denominator := new(big.Int).Lsh(big.NewInt(1), offset)
	div := new(big.Int).Div(prod, denominator)
	if roundUp {
		mod := new(big.Int).Mod(prod, denominator)
		if mod.Sign() > 0 {
			div.Add(div, big.NewInt(1))
		}
	}
	return div
}

// ShlDiv performs: (x << offset) / y, optionally rounding up if there is a remainder.
func ShlDiv(x, y *big.Int, offset uint, roundUp bool) *big.Int {
	num := new(big.Int).Lsh(x, offset)
	div := new(big.Int).Div(num, y)
	if roundUp {
		mod := new(big.Int).Mod(num, y)
		if mod.Sign() > 0 {
			div.Add(div, big.NewInt(1))
		}
	}
	return div
}

// GetAmountIn calculates required input tokens given output tokens.
func GetAmountIn(amountOut, price *big.Int, swapForY bool, roundUp bool) *big.Int {
	if swapForY {
		return ShlDiv(amountOut, price, 64, roundUp)
	}
	return MulShr(amountOut, price, 64, roundUp)
}

// GetAmountOut calculates received output tokens given input tokens.
func GetAmountOut(price, inAmount *big.Int, swapForY bool, roundUp bool) *big.Int {
	if swapForY {
		return MulShr(inAmount, price, 64, roundUp)
	}
	return ShlDiv(inAmount, price, 64, roundUp)
}

// BinIdToPriceDecimal calculates the price of a bin in lamport relative price with arbitrary precision.
// Price = (1 + binStep/10000) ^ binId
func BinIdToPriceDecimal(binId int32, binStep uint16) decimal.Decimal {
	binStepDec := decimal.NewFromInt(int64(binStep))
	bpsDec := decimal.NewFromInt(10000)
	base := decimal.NewFromInt(1).Add(binStepDec.Div(bpsDec))
	return base.Pow(decimal.NewFromInt(int64(binId)))
}

// BinIdToPrice calculates the price of a bin in lamport relative price.
func BinIdToPrice(binId int32, binStep uint16) float64 {
	val := BinIdToPriceDecimal(binId, binStep)
	f, _ := val.Float64()
	return f
}

// PriceToBinIdDecimal calculates the bin ID closest to a given lamport price represented as decimal.
func PriceToBinIdDecimal(price decimal.Decimal, binStep uint16) int32 {
	priceF, _ := price.Float64()
	base := 1.0 + float64(binStep)/10000.0
	binId := math.Log(priceF) / math.Log(base)
	return int32(math.Round(binId))
}

// PriceToBinId calculates the bin ID closest to a given lamport price.
func PriceToBinId(price float64, binStep uint16) int32 {
	return PriceToBinIdDecimal(decimal.NewFromFloat(price), binStep)
}

// BinIdToTokenPriceDecimal calculates the user-facing price (token Y per token X) adjusting for decimals with arbitrary precision.
func BinIdToTokenPriceDecimal(binId int32, binStep uint16, decimalsX int32, decimalsY int32) decimal.Decimal {
	pricePerLamport := BinIdToPriceDecimal(binId, binStep)
	diff := int64(decimalsX - decimalsY)
	multiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(diff))
	return pricePerLamport.Mul(multiplier)
}

// BinIdToTokenPrice calculates the user-facing price (token Y per token X) adjusting for decimals.
func BinIdToTokenPrice(binId int32, binStep uint16, decimalsX int32, decimalsY int32) float64 {
	val := BinIdToTokenPriceDecimal(binId, binStep, decimalsX, decimalsY)
	f, _ := val.Float64()
	return f
}

// TokenPriceToBinIdDecimal calculates the bin ID for a given user-facing token price represented as decimal.
func TokenPriceToBinIdDecimal(tokenPrice decimal.Decimal, binStep uint16, decimalsX int32, decimalsY int32) int32 {
	diff := int64(decimalsX - decimalsY)
	multiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(diff))
	pricePerLamport := tokenPrice.Div(multiplier)
	return PriceToBinIdDecimal(pricePerLamport, binStep)
}

// TokenPriceToBinId calculates the bin ID for a given user-facing token price.
func TokenPriceToBinId(tokenPrice float64, binStep uint16, decimalsX int32, decimalsY int32) int32 {
	return TokenPriceToBinIdDecimal(decimal.NewFromFloat(tokenPrice), binStep, decimalsX, decimalsY)
}

// GetBaseFee calculates the base swap fee numerator.
func GetBaseFee(binStep uint16, sParam StaticParameters) *big.Int {
	baseFactor := big.NewInt(int64(sParam.BaseFactor))
	step := big.NewInt(int64(binStep))
	ten := big.NewInt(10)
	power := new(big.Int).Exp(ten, big.NewInt(int64(sParam.BaseFeePowerFactor)), nil)

	res := new(big.Int).Mul(baseFactor, step)
	res.Mul(res, ten)
	res.Mul(res, power)
	return res
}

// GetVariableFee calculates the variable fee numerator.
func GetVariableFee(binStep uint16, sParam StaticParameters, vParam VariableParameters) *big.Int {
	if sParam.VariableFeeControl > 0 {
		vfa := big.NewInt(int64(vParam.VolatilityAccumulator))
		step := big.NewInt(int64(binStep))

		squareVfaBin := new(big.Int).Mul(vfa, step)
		squareVfaBin.Mul(squareVfaBin, squareVfaBin)

		vFee := new(big.Int).Mul(big.NewInt(int64(sParam.VariableFeeControl)), squareVfaBin)

		vFee.Add(vFee, big.NewInt(99_999_999_999))
		vFee.Div(vFee, big.NewInt(100_000_000_000))
		return vFee
	}
	return big.NewInt(0)
}

// GetTotalFee calculates the total fee numerator.
func GetTotalFee(binStep uint16, sParam StaticParameters, vParam VariableParameters) *big.Int {
	baseFee := GetBaseFee(binStep, sParam)
	varFee := GetVariableFee(binStep, sParam, vParam)
	totalFee := new(big.Int).Add(baseFee, varFee)
	if totalFee.Cmp(MaxFeeRate) > 0 {
		return MaxFeeRate
	}
	return totalFee
}

// FeeAmountResult contains the resulting amount and fee.
type FeeAmountResult struct {
	Amount *big.Int
	Fee    *big.Int
}

// GetExcludedFeeAmount returns the excluded fee amount and fee.
func GetExcludedFeeAmount(includedFeeAmount, tradeFeeNumerator *big.Int) FeeAmountResult {
	feePrecMinus1 := new(big.Int).Sub(FeePrecision, big.NewInt(1))
	num := new(big.Int).Mul(includedFeeAmount, tradeFeeNumerator)
	num.Add(num, feePrecMinus1)
	tradingFee := new(big.Int).Div(num, FeePrecision)

	excludedFeeAmount := new(big.Int).Sub(includedFeeAmount, tradingFee)
	return FeeAmountResult{
		Amount: excludedFeeAmount,
		Fee:    tradingFee,
	}
}

// GetIncludedFeeAmount returns the included fee amount and fee.
func GetIncludedFeeAmount(excludedFeeAmount, tradeFeeNumerator *big.Int) FeeAmountResult {
	denominator := new(big.Int).Sub(FeePrecision, tradeFeeNumerator)

	denomMinus1 := new(big.Int).Sub(denominator, big.NewInt(1))
	num := new(big.Int).Mul(excludedFeeAmount, FeePrecision)
	num.Add(num, denomMinus1)
	includedFeeAmount := new(big.Int).Div(num, denominator)

	fee := new(big.Int).Sub(includedFeeAmount, excludedFeeAmount)
	return FeeAmountResult{
		Amount: includedFeeAmount,
		Fee:    fee,
	}
}

// SplitFeeResult contains split fees.
type SplitFeeResult struct {
	Fee         *big.Int
	ProtocolFee *big.Int
}

// SplitFee splits trading fees between user and protocol.
func SplitFee(tradingFee, protocolShare, mmAmountIn, totalAmountIn *big.Int) SplitFeeResult {
	if totalAmountIn.Sign() == 0 {
		return SplitFeeResult{
			Fee:         big.NewInt(0),
			ProtocolFee: big.NewInt(0),
		}
	}

	totalMinus1 := new(big.Int).Sub(totalAmountIn, big.NewInt(1))
	num := new(big.Int).Mul(tradingFee, mmAmountIn)
	num.Add(num, totalMinus1)
	mmFee := new(big.Int).Div(num, totalAmountIn)

	totalLoFee := new(big.Int).Sub(tradingFee, mmFee)

	loFee := new(big.Int).Mul(totalLoFee, LimitOrderFeeShare)
	loFee.Div(loFee, BasisPointMaxBig)

	loProtocolFee := new(big.Int).Sub(totalLoFee, loFee)

	mmProtocolFee := new(big.Int).Mul(mmFee, protocolShare)
	mmProtocolFee.Div(mmProtocolFee, BasisPointMaxBig)

	totalProtocolFee := new(big.Int).Add(loProtocolFee, mmProtocolFee)
	totalUserFee := new(big.Int).Sub(tradingFee, totalProtocolFee)

	return SplitFeeResult{
		Fee:         totalUserFee,
		ProtocolFee: totalProtocolFee,
	}
}

// FeeMode contains fee calculation modes.
type FeeMode struct {
	FeeOnInput  bool
	FeeOnTokenX bool
}

// GetFeeMode returns the fee mode based on parameters.
func GetFeeMode(collectFeeMode uint8, swapForY bool) FeeMode {
	var feeOnInput bool
	var feeOnTokenX bool

	switch collectFeeMode {
	case 0: // InputOnly
		feeOnInput = true
		feeOnTokenX = swapForY
	case 1: // OnlyY
		feeOnInput = !swapForY
		feeOnTokenX = false
	default:
		// Default to InputOnly if unrecognized
		feeOnInput = true
		feeOnTokenX = swapForY
	}

	return FeeMode{
		FeeOnInput:  feeOnInput,
		FeeOnTokenX: feeOnTokenX,
	}
}

// GetLimitOrderLiquidity extracts limit order liquidity from a bin.
func GetLimitOrderLiquidity(bin Bin, supportLimitOrder bool) (orderAmountX, orderAmountY *big.Int) {
	if !supportLimitOrder {
		return big.NewInt(0), big.NewInt(0)
	}

	totalOrderAmount := new(big.Int).Add(
		big.NewInt(0).SetUint64(bin.OpenOrderAmount),
		big.NewInt(0).SetUint64(bin.ProcessedOrderRemainingAmount),
	)

	isAskSide := bin.LimitOrderAskSide != 0
	if isAskSide {
		return totalOrderAmount, big.NewInt(0)
	}
	return big.NewInt(0), totalOrderAmount
}

// GetBinMaxAmountOut calculates the max output amount possible from a bin.
func GetBinMaxAmountOut(bin Bin, swapForY bool, supportLimitOrder bool) *big.Int {
	orderX, orderY := GetLimitOrderLiquidity(bin, supportLimitOrder)
	if swapForY {
		return new(big.Int).Add(big.NewInt(0).SetUint64(bin.AmountY), orderY)
	}
	return new(big.Int).Add(big.NewInt(0).SetUint64(bin.AmountX), orderX)
}

type FillAmountResult struct {
	AmountIn    *big.Int
	AmountLeft  *big.Int
	OutAmount   *big.Int
	MmAmountIn  *big.Int
	MmAmountOut *big.Int
}

func calculateExactInFillAmount(bin Bin, amount, maxAmountOut *big.Int, swapForY bool) FillAmountResult {
	binPrice := U128ToBigInt(bin.Price)
	maxAmountIn := GetAmountIn(maxAmountOut, binPrice, swapForY, true)

	if amount.Cmp(maxAmountIn) >= 0 {
		return FillAmountResult{
			AmountIn:    maxAmountIn,
			AmountLeft:  new(big.Int).Sub(amount, maxAmountIn),
			OutAmount:   maxAmountOut,
			MmAmountIn:  maxAmountIn,
			MmAmountOut: maxAmountOut,
		}
	}

	outAmount := GetAmountOut(binPrice, amount, swapForY, false)
	return FillAmountResult{
		AmountIn:    amount,
		AmountLeft:  big.NewInt(0),
		OutAmount:   outAmount,
		MmAmountIn:  amount,
		MmAmountOut: outAmount,
	}
}

func getLimitOrderAmountsBySwapDirection(bin Bin, swapForY bool) (openOrderAmount, processedOrderRemainingAmount *big.Int) {
	isAskSide := bin.LimitOrderAskSide != 0
	if (swapForY && !isAskSide) || (!swapForY && isAskSide) {
		return big.NewInt(0).SetUint64(bin.OpenOrderAmount), big.NewInt(0).SetUint64(bin.ProcessedOrderRemainingAmount)
	}
	return big.NewInt(0), big.NewInt(0)
}

func getExactInFillAmountResult(bin Bin, amountIn *big.Int, swapForY bool, supportLimitOrder bool) FillAmountResult {
	mmAmount := big.NewInt(0)
	if swapForY {
		mmAmount.SetUint64(bin.AmountY)
	} else {
		mmAmount.SetUint64(bin.AmountX)
	}

	mmFill := calculateExactInFillAmount(bin, amountIn, mmAmount, swapForY)

	if !supportLimitOrder {
		return FillAmountResult{
			AmountIn:    mmFill.AmountIn,
			AmountLeft:  mmFill.AmountLeft,
			OutAmount:   mmFill.OutAmount,
			MmAmountIn:  mmFill.AmountIn,
			MmAmountOut: mmFill.OutAmount,
		}
	}

	amountLeftAfterMM := mmFill.AmountLeft
	processedOrderAmountIn := big.NewInt(0)
	processedOrderAmountOut := big.NewInt(0)
	openOrderAmountIn := big.NewInt(0)
	openOrderAmountOut := big.NewInt(0)

	if amountLeftAfterMM.Sign() > 0 {
		openOrderAmount, processedOrderRemainingAmount := getLimitOrderAmountsBySwapDirection(bin, swapForY)

		remainingOrderFill := calculateExactInFillAmount(bin, amountLeftAfterMM, processedOrderRemainingAmount, swapForY)
		processedOrderAmountIn = remainingOrderFill.AmountIn
		processedOrderAmountOut = remainingOrderFill.OutAmount

		if remainingOrderFill.AmountLeft.Sign() > 0 {
			openOrderFill := calculateExactInFillAmount(bin, remainingOrderFill.AmountLeft, openOrderAmount, swapForY)
			openOrderAmountIn = openOrderFill.AmountIn
			openOrderAmountOut = openOrderFill.OutAmount
		}
	}

	totalAmountIn := new(big.Int).Add(mmFill.AmountIn, processedOrderAmountIn)
	totalAmountIn.Add(totalAmountIn, openOrderAmountIn)

	totalAmountOut := new(big.Int).Add(mmFill.OutAmount, processedOrderAmountOut)
	totalAmountOut.Add(totalAmountOut, openOrderAmountOut)

	return FillAmountResult{
		AmountIn:    totalAmountIn,
		AmountLeft:  new(big.Int).Sub(amountIn, totalAmountIn),
		OutAmount:   totalAmountOut,
		MmAmountIn:  mmFill.AmountIn,
		MmAmountOut: mmFill.OutAmount,
	}
}

type SwapAtBinResult struct {
	AmountIn    *big.Int
	AmountOut   *big.Int
	Fee         *big.Int
	ProtocolFee *big.Int
}

// SwapExactInQuoteAtBin simulates swap at a specific bin.
func SwapExactInQuoteAtBin(
	bin Bin,
	binStep uint16,
	sParam StaticParameters,
	vParam VariableParameters,
	inAmount *big.Int,
	swapForY bool,
	supportLimitOrder bool,
	feeOnInput bool,
) SwapAtBinResult {
	tradingFee := big.NewInt(0)
	excludedFeeAmountIn := inAmount

	tradeFeeNumerator := GetTotalFee(binStep, sParam, vParam)

	if feeOnInput {
		res := GetExcludedFeeAmount(inAmount, tradeFeeNumerator)
		tradingFee = res.Fee
		excludedFeeAmountIn = res.Amount
	}

	fillAmountResult := getExactInFillAmountResult(bin, excludedFeeAmountIn, swapForY, supportLimitOrder)

	amountLeft := fillAmountResult.AmountLeft
	outAmount := fillAmountResult.OutAmount

	includedFeeAmountIn := inAmount

	if amountLeft.Sign() > 0 {
		excludedFeeAmountIn = new(big.Int).Sub(excludedFeeAmountIn, amountLeft)

		if feeOnInput {
			res := GetIncludedFeeAmount(excludedFeeAmountIn, tradeFeeNumerator)
			tradingFee = res.Fee
			includedFeeAmountIn = res.Amount
		} else {
			includedFeeAmountIn = excludedFeeAmountIn
		}
	}

	excludedFeeAmountOut := outAmount

	if !feeOnInput {
		res := GetExcludedFeeAmount(outAmount, tradeFeeNumerator)
		tradingFee = res.Fee
		excludedFeeAmountOut = res.Amount
	}

	splitRes := SplitFee(tradingFee, big.NewInt(int64(sParam.ProtocolShare)), fillAmountResult.MmAmountIn, fillAmountResult.AmountIn)

	return SwapAtBinResult{
		AmountIn:    includedFeeAmountIn,
		AmountOut:   excludedFeeAmountOut,
		Fee:         splitRes.Fee,
		ProtocolFee: splitRes.ProtocolFee,
	}
}

// UpdateVolatilityAccumulator updates the volatility accumulator.
// Matches TS SDK index.ts:9289-9302 exactly.
func UpdateVolatilityAccumulator(vParam *VariableParameters, sParam StaticParameters, activeId int32) {
	deltaId := activeId - vParam.IndexReference
	if deltaId < 0 {
		deltaId = -deltaId
	}
	newVolatilityAccumulator := vParam.VolatilityReference + uint32(deltaId)*10000
	if newVolatilityAccumulator > sParam.MaxVolatilityAccumulator {
		vParam.VolatilityAccumulator = sParam.MaxVolatilityAccumulator
	} else {
		vParam.VolatilityAccumulator = newVolatilityAccumulator
	}
}

// UpdateReference updates volatility parameters based on reference timestamp.
// Matches TS SDK index.ts:9304-9325 exactly.
// NOTE: TS does NOT set VolatilityAccumulator or LastUpdateTimestamp here.
func UpdateReference(activeId int32, vParam *VariableParameters, sParam StaticParameters, currentTimestamp int64) {
	elapsed := currentTimestamp - vParam.LastUpdateTimestamp
	if elapsed >= int64(sParam.FilterPeriod) {
		vParam.IndexReference = activeId
		if elapsed < int64(sParam.DecayPeriod) {
			decayedVolatilityReference := (uint64(vParam.VolatilityAccumulator) * uint64(sParam.ReductionFactor)) / 10000
			vParam.VolatilityReference = uint32(decayedVolatilityReference)
		} else {
			vParam.VolatilityReference = 0
		}
	}
}

// ComputeFeeFromAmount computes the fee from a fee-inclusive amount.
// Matches TS SDK fee.ts:60-71 exactly.
func ComputeFeeFromAmount(binStep uint16, sParam StaticParameters, vParam VariableParameters, inAmountWithFees *big.Int) *big.Int {
	totalFee := GetTotalFee(binStep, sParam, vParam)
	feePrecMinus1 := new(big.Int).Sub(FeePrecision, big.NewInt(1))
	num := new(big.Int).Mul(inAmountWithFees, totalFee)
	num.Add(num, feePrecMinus1)
	return new(big.Int).Div(num, FeePrecision)
}
