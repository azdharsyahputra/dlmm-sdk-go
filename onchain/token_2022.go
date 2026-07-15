package onchain

import (
	"math/big"
)

// MAX_FEE_BASIS_POINTS is 10000 (100%).
var MaxFeeBasisPoints = big.NewInt(10000)

// TransferFee contains the Token-2022 transfer fee configuration for a token.
type TransferFee struct {
	TransferFeeBasisPoints uint16
	MaximumFee             uint64
}

// CalculateFee computes the transfer fee for a given amount.
func CalculateFee(transferFee *TransferFee, amount *big.Int) *big.Int {
	if transferFee == nil || transferFee.TransferFeeBasisPoints == 0 {
		return big.NewInt(0)
	}
	if amount.Sign() == 0 {
		return big.NewInt(0)
	}
	if transferFee.TransferFeeBasisPoints == 10000 {
		return new(big.Int).SetUint64(transferFee.MaximumFee)
	}

	bps := big.NewInt(int64(transferFee.TransferFeeBasisPoints))
	fee := new(big.Int).Mul(amount, bps)
	fee.Add(fee, new(big.Int).Sub(MaxFeeBasisPoints, big.NewInt(1)))
	fee.Div(fee, MaxFeeBasisPoints)

	maxFee := new(big.Int).SetUint64(transferFee.MaximumFee)
	if fee.Cmp(maxFee) > 0 {
		return maxFee
	}
	return fee
}

// CalculatePreFeeAmount computes the pre-fee amount given a post-fee amount (inverse fee calculation).
func CalculatePreFeeAmount(transferFee *TransferFee, postFeeAmount *big.Int) *big.Int {
	if transferFee == nil || transferFee.TransferFeeBasisPoints == 0 {
		return new(big.Int).Set(postFeeAmount)
	}
	if postFeeAmount.Sign() == 0 {
		return big.NewInt(0)
	}

	maxFee := new(big.Int).SetUint64(transferFee.MaximumFee)

	if transferFee.TransferFeeBasisPoints == 10000 {
		return new(big.Int).Add(postFeeAmount, maxFee)
	}

	numerator := new(big.Int).Mul(postFeeAmount, MaxFeeBasisPoints)
	denominator := new(big.Int).Sub(MaxFeeBasisPoints, big.NewInt(int64(transferFee.TransferFeeBasisPoints)))

	rawPreFeeAmount := new(big.Int).Add(numerator, denominator)
	rawPreFeeAmount.Sub(rawPreFeeAmount, big.NewInt(1))
	rawPreFeeAmount.Div(rawPreFeeAmount, denominator)

	diff := new(big.Int).Sub(rawPreFeeAmount, postFeeAmount)
	if diff.Cmp(maxFee) >= 0 {
		return new(big.Int).Add(postFeeAmount, maxFee)
	}
	return rawPreFeeAmount
}

// CalculateTransferFeeIncludedAmount calculates the total amount needed (including fee) to yield the specified excluded amount.
func CalculateTransferFeeIncludedAmount(transferFee *TransferFee, transferFeeExcludedAmount *big.Int) (*big.Int, *big.Int) {
	if transferFee == nil || transferFeeExcludedAmount.Sign() == 0 {
		return new(big.Int).Set(transferFeeExcludedAmount), big.NewInt(0)
	}

	var fee *big.Int
	if transferFee.TransferFeeBasisPoints == 10000 {
		fee = new(big.Int).SetUint64(transferFee.MaximumFee)
	} else {
		preFeeAmount := CalculatePreFeeAmount(transferFee, transferFeeExcludedAmount)
		fee = CalculateFee(transferFee, preFeeAmount)
	}

	includedAmount := new(big.Int).Add(transferFeeExcludedAmount, fee)
	return includedAmount, fee
}

// CalculateTransferFeeExcludedAmount calculates the final amount (excluding fee) yielded from the specified included amount.
func CalculateTransferFeeExcludedAmount(transferFee *TransferFee, transferFeeIncludedAmount *big.Int) (*big.Int, *big.Int) {
	if transferFee == nil || transferFeeIncludedAmount.Sign() == 0 {
		return new(big.Int).Set(transferFeeIncludedAmount), big.NewInt(0)
	}

	fee := CalculateFee(transferFee, transferFeeIncludedAmount)
	excludedAmount := new(big.Int).Sub(transferFeeIncludedAmount, fee)
	return excludedAmount, fee
}
