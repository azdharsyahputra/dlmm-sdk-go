package onchain_test

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func ExampleComputeSwapQuote() {
	// 1. Connect to Solana RPC
	rpcClient := rpc.New(rpc.MainNetBeta_RPC)

	// 2. Fetch LbPair state
	poolAddress := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")
	pair, err := onchain.FetchLbPair(context.Background(), rpcClient, poolAddress)
	if err != nil {
		log.Fatalf("Failed to fetch LbPair: %v", err)
	}

	// 3. Fetch surrounding BinArrays
	binArrayIndex := onchain.BinIdToBinArrayIndex(pair.ActiveId)
	binArrays, err := onchain.FetchBinArrays(context.Background(), rpcClient, poolAddress, binArrayIndex, 3)
	if err != nil {
		log.Fatalf("Failed to fetch BinArrays: %v", err)
	}

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

func ExampleBuildSwapInstruction() {
	// Assume 'quote' was obtained from ComputeSwapQuote()
	// and 'pair' from FetchLbPair().
	var quote *onchain.SwapQuote
	var pair *onchain.LbPair
	userWallet := solana.MustPublicKeyFromBase58("YourWalletAddressHere")
	userTokenXAccount := solana.MustPublicKeyFromBase58("YourTokenXAccount")
	userTokenYAccount := solana.MustPublicKeyFromBase58("YourTokenYAccount")

	// The pool pubkey (LB Pair)
	poolPubkey := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")

	inst, err := onchain.BuildSwapInstruction(
		poolPubkey,
		pair,
		userWallet,
		userTokenXAccount,
		userTokenYAccount,
		quote,
	)
	if err != nil {
		log.Fatalf("Failed to build instruction: %v", err)
	}

	fmt.Printf("Instruction ProgramID: %s\n", inst.ProgramID().String())
}
