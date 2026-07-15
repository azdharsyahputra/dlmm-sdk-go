package onchain

import (
	"encoding/binary"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

// ProgramID is the Meteora DLMM on-chain program address.
var ProgramID = solana.MustPublicKeyFromBase58("LBUZKhRxPF3XUpBCjp4YzTKgLccjZhTSDM9YuVaPwxo")

// MaxBinArraySize is the number of bins per BinArray account.
const MaxBinArraySize int64 = 70

// DeriveBinArray derives the PDA for a BinArray at the given index.
// Uses two's complement for negative indices, matching the TS SDK exactly.
func DeriveBinArray(lbPair solana.PublicKey, index int64) (solana.PublicKey, uint8, error) {
	indexBytes := make([]byte, 8)
	if index < 0 {
		// Two's complement for negative i64
		twos := new(big.Int).SetInt64(index)
		twos.Add(twos, new(big.Int).Lsh(big.NewInt(1), 64))
		b := twos.Bytes()
		// Convert big-endian to little-endian
		for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
			b[i], b[j] = b[j], b[i]
		}
		copy(indexBytes, b)
	} else {
		binary.LittleEndian.PutUint64(indexBytes, uint64(index))
	}
	seeds := [][]byte{
		[]byte("bin_array"),
		lbPair.Bytes(),
		indexBytes,
	}
	pubkey, bump, err := solana.FindProgramAddress(seeds, ProgramID)
	return pubkey, bump, err
}

// DeriveReserve derives the PDA for a vault token reserve.
// Seed order: [lbPair, token] — matching the TS SDK (derive.ts:201-203).
func DeriveReserve(token solana.PublicKey, lbPair solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		lbPair.Bytes(),
		token.Bytes(),
	}
	pubkey, bump, err := solana.FindProgramAddress(seeds, ProgramID)
	return pubkey, bump, err
}

// DeriveBinArrayBitmapExtension derives the PDA for the bitmap extension account.
func DeriveBinArrayBitmapExtension(lbPair solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("bitmap"),
		lbPair.Bytes(),
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveOracle derives the PDA for the oracle account.
func DeriveOracle(lbPair solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("oracle"),
		lbPair.Bytes(),
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DerivePosition derives the PDA for a position account.
// Uses two's complement for negative lowerBinID.
func DerivePosition(lbPair solana.PublicKey, base solana.PublicKey, lowerBinID int32, width int32) (solana.PublicKey, uint8, error) {
	lowerBinIDBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lowerBinIDBytes, uint32(lowerBinID)) // Handles two's complement for negative automatically

	widthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(widthBytes, uint32(width))

	seeds := [][]byte{
		[]byte("position"),
		lbPair.Bytes(),
		base.Bytes(),
		lowerBinIDBytes,
		widthBytes,
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveLBPair derives the PDA for an LBPair from token mints and bin step.
func DeriveLBPair(tokenX solana.PublicKey, tokenY solana.PublicKey, binStep uint16) (solana.PublicKey, uint8, error) {
	minKey, maxKey := sortTokenMints(tokenX, tokenY)

	binStepBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(binStepBytes, binStep)

	seeds := [][]byte{
		minKey.Bytes(),
		maxKey.Bytes(),
		binStepBytes,
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveLBPair2 derives the PDA for an LBPair v2 from token mints, bin step, and base factor.
func DeriveLBPair2(tokenX solana.PublicKey, tokenY solana.PublicKey, binStep uint16, baseFactor uint16) (solana.PublicKey, uint8, error) {
	minKey, maxKey := sortTokenMints(tokenX, tokenY)

	binStepBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(binStepBytes, binStep)

	baseFactorBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(baseFactorBytes, baseFactor)

	seeds := [][]byte{
		minKey.Bytes(),
		maxKey.Bytes(),
		binStepBytes,
		baseFactorBytes,
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveRewardVault derives the PDA for a reward vault.
func DeriveRewardVault(lbPair solana.PublicKey, rewardIndex uint64) (solana.PublicKey, uint8, error) {
	rewardIndexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(rewardIndexBytes, rewardIndex)

	seeds := [][]byte{
		lbPair.Bytes(),
		rewardIndexBytes,
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DerivePresetParameter derives the PDA for a preset parameter.
func DerivePresetParameter(binStep uint16) (solana.PublicKey, uint8, error) {
	binStepBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(binStepBytes, binStep)

	seeds := [][]byte{
		[]byte("preset_parameter"),
		binStepBytes,
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// sortTokenMints returns tokens sorted by their byte representation (ascending).
func sortTokenMints(tokenX solana.PublicKey, tokenY solana.PublicKey) (solana.PublicKey, solana.PublicKey) {
	xBytes := tokenX.Bytes()
	yBytes := tokenY.Bytes()
	for i := 0; i < 32; i++ {
		if xBytes[i] > yBytes[i] {
			return tokenY, tokenX
		} else if xBytes[i] < yBytes[i] {
			return tokenX, tokenY
		}
	}
	return tokenX, tokenY
}

// BinIDToBinArrayIndex returns the index of the BinArray containing the binID.
// Equivalent to TypeScript binID.divmod(70) behavior with floor division for negatives.
func BinIDToBinArrayIndex(binID int32) int64 {
	idx := int64(binID) / MaxBinArraySize
	mod := int64(binID) % MaxBinArraySize
	if binID < 0 && mod != 0 {
		return idx - 1
	}
	return idx
}

// GetBinArrayLowerUpperBinID returns the lower and upper bin ID limits of a BinArray index.
func GetBinArrayLowerUpperBinID(binArrayIndex int64) (int32, int32) {
	lowerBinID := int32(binArrayIndex * MaxBinArraySize)
	upperBinID := lowerBinID + int32(MaxBinArraySize) - 1
	return lowerBinID, upperBinID
}

// ilmBase is the base public key used for customizable permissionless LBPairs.
var ilmBase = solana.MustPublicKeyFromBase58("MFGQxwAmB91SwuYX36okv2Qmdc9aMuHTwWGUrp4AtB1")

// DerivePresetParameterWithIndex derives the PDA for a preset parameter with index.
func DerivePresetParameterWithIndex(index uint16) (solana.PublicKey, uint8, error) {
	indexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(indexBytes, index)

	seeds := [][]byte{
		[]byte("preset_parameter2"),
		indexBytes,
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveLBPairWithPresetParamWithIndexKey derives the PDA for an LBPair using a preset parameter key.
func DeriveLBPairWithPresetParamWithIndexKey(presetParameterKey solana.PublicKey, tokenX solana.PublicKey, tokenY solana.PublicKey) (solana.PublicKey, uint8, error) {
	minKey, maxKey := sortTokenMints(tokenX, tokenY)
	seeds := [][]byte{
		presetParameterKey.Bytes(),
		minKey.Bytes(),
		maxKey.Bytes(),
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveCustomizablePermissionlessLBPair derives the PDA for a customizable permissionless LBPair.
func DeriveCustomizablePermissionlessLBPair(tokenX solana.PublicKey, tokenY solana.PublicKey) (solana.PublicKey, uint8, error) {
	minKey, maxKey := sortTokenMints(tokenX, tokenY)
	seeds := [][]byte{
		ilmBase.Bytes(),
		minKey.Bytes(),
		maxKey.Bytes(),
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DerivePermissionLBPair derives the PDA for a permissioned LBPair.
func DerivePermissionLBPair(baseKey solana.PublicKey, tokenX solana.PublicKey, tokenY solana.PublicKey, binStep uint16) (solana.PublicKey, uint8, error) {
	minKey, maxKey := sortTokenMints(tokenX, tokenY)
	binStepBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(binStepBytes, binStep)

	seeds := [][]byte{
		baseKey.Bytes(),
		minKey.Bytes(),
		maxKey.Bytes(),
		binStepBytes,
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveTokenBadge derives the PDA for a token badge.
func DeriveTokenBadge(mint solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("token_badge"),
		mint.Bytes(),
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// DeriveOperator derives the PDA for an operator.
func DeriveOperator(whitelistedSigner solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("operator"),
		whitelistedSigner.Bytes(),
	}
	pubkey, bump, err := solana.FindProgramAddress(seeds, ProgramID)
	return pubkey, bump, err
}
