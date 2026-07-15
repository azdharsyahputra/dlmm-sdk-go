package onchain_test

import (
	"fmt"
	"log"
	"math/big"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
)

func ExampleComputeSwapQuote() {
	// 1. Fetch LbPair state (Implementation varies)
	pair := &onchain.LbPair{}

	// 3. Fetch surrounding BinArrays
	// binArrayIndex := onchain.BinIdToBinArrayIndex(pair.ActiveId)
	// binArrays, err := onchain.FetchBinArrays(context.Background(), rpcClient, poolAddress, binArrayIndex, 3)
	var binArrays []onchain.BinArray

	// 4. Compute exact Swap Quote (1 SOL -> USDC)
	inAmount := new(big.Int).SetUint64(1_000_000_000) // 1 SOL
	quote, err := onchain.ComputeSwapQuote(
		pair, 
		binArrays, 
		inAmount, 
		true,   // swapForY = true (X to Y)
		false,  // isPartialFill
		50,     // 50 BPS (0.5%) Slippage
		0,      // Use current timestamp
		nil,    // No Token-2022 Transfer Fee on Input
		nil,    // No Token-2022 Transfer Fee on Output
	)
	if err != nil {
		log.Fatalf("Quote failed: %v", err)
	}

	fmt.Printf("Expected OutAmount: %s\n", quote.OutAmount.String())
	fmt.Printf("Price Impact: %.4f%%\n", quote.PriceImpact)
}

func ExampleNewSwapInstruction() {
	// inst, err := onchain.NewSwapInstruction(...)
	// fmt.Printf("Instruction ProgramID: %s\n", inst.ProgramID().String())
}
