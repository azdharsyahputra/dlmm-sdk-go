package onchain_test

import (
	"fmt"
	"log"
	"math/big"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
)

func ExampleDecodeLBPair() {
	var pair onchain.LBPair
	err := onchain.DecodeLBPair([]byte{1, 2, 3}, &pair)
	if err != nil {
		log.Fatalf("Failed to fetch LBPair: %v", err)
	}

	fmt.Printf("Fetched Pool! Active Bin ID: %d\n", pair.ActiveID)
}

func ExampleComputeSwapQuote() {
	// 1. Setup mock or fetched data
	pair := &onchain.LBPair{
		ActiveID: -100,
		BinStep:  10,
	}
	binArrays := []*onchain.BinArray{} // Need to fetch BinArrays based on ActiveID

	inAmount := new(big.Int).SetUint64(1_000_000) // 1 Token
	swapForY := true                              // Swap X -> Y

	// 2. Compute quote
	quote, err := onchain.ComputeSwapQuote(
		pair,
		binArrays,
		inAmount,
		swapForY,
		false,
		50, // 50 bps slippage (0.5%)
		0,
		nil,
		nil,
	)
	if err != nil {
		log.Printf("Failed to compute quote (expected without real bins): %v", err)
		return
	}

	fmt.Printf("Expected Out Amount: %s\n", quote.OutAmount.String())
}

func ExampleNewSwapInstruction() {
	poolAddress := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")
	userWallet := solana.MustPublicKeyFromBase58("2wT8Yq49kHgDzFolAW5dK3a7R228G7k4A2m8j1x6zZ1x")

	// These parameters would normally come from your fetched pool state and ComputeSwapQuote
	params := onchain.SwapParams{
		LBPair:        poolAddress,
		ReserveX:      solana.NewWallet().PublicKey(),
		ReserveY:      solana.NewWallet().PublicKey(),
		UserTokenIn:   userWallet,
		UserTokenOut:  userWallet,
		TokenXMint:    solana.NewWallet().PublicKey(),
		TokenYMint:    solana.NewWallet().PublicKey(),
		Oracle:        solana.NewWallet().PublicKey(),
		User:          userWallet,
		TokenXProgram: solana.TokenProgramID,
		TokenYProgram: solana.TokenProgramID,
		BinArrays:     []solana.PublicKey{solana.NewWallet().PublicKey()},
		AmountIn:      1_000_000,
		MinAmountOut:  990_000,
	}

	inst, err := onchain.NewSwapInstruction(params)
	if err != nil {
		log.Fatalf("Failed to build swap instruction: %v", err)
	}

	fmt.Printf("Instruction Program ID: %s\n", inst.ProgramID().String())
}

func ExampleBinIDToTokenPrice() {
	// SOL-USDC Pool typically has BinStep=4, decimals X=9 (SOL), decimals Y=6 (USDC)
	binID := int32(-6368)
	binStep := uint16(4)
	decimalsX := int32(9)
	decimalsY := int32(6)

	// Calculate float64 price
	price := onchain.BinIDToTokenPrice(binID, binStep, decimalsX, decimalsY)

	fmt.Printf("Price of SOL at Bin %d is approx %.2f USDC\n", binID, price)
}
