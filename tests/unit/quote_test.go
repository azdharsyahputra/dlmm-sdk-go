package unit_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TestSwapQuoteSOLUSDC(t *testing.T) {
	c := onchain.NewClient(rpc.MainNetBeta_RPC)
	ctx := context.Background()

	// SOL-USDC pool address on mainnet
	poolAddr := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")

	// 1. Fetch LBPair
	pair, err := c.GetLBPair(ctx, poolAddr)
	if err != nil {
		t.Fatalf("Failed to fetch LBPair: %v", err)
	}

	t.Logf("SOL-USDC Pool Active Bin ID: %d", pair.ActiveID)

	// 2. Fetch the active BinArray and adjacent ones
	activeIdx := onchain.BinIDToBinArrayIndex(pair.ActiveID)
	t.Logf("Active BinArray Index: %d", activeIdx)

	// We will fetch active index, active-1, and active+1 to cover price movement in both directions
	indices := []int64{activeIdx - 1, activeIdx, activeIdx + 1}
	var binArrays []*onchain.BinArray

	for _, idx := range indices {
		ba, err := c.GetBinArrayByIndex(ctx, poolAddr, idx)
		if err == nil {
			binArrays = append(binArrays, ba)
			t.Logf("Fetched BinArray index %d", idx)
		} else {
			t.Logf("Warning: Could not fetch BinArray index %d (might not be initialized): %v", idx, err)
		}
	}

	// 3. Simulating Swap Exact-In: 1 SOL -> USDC (swapForY = true)
	// Swap exact in: 1 SOL -> USDC
	inAmount := new(big.Int).SetUint64(1_000_000_000)
	quoteY, err := onchain.ComputeSwapQuote(pair, binArrays, inAmount, true, false, 50, 0, nil, nil)
	if err != nil {
		t.Fatalf("Failed to compute swap quote Y: %v", err)
	}

	t.Logf("--- SOL -> USDC Swap Quote (1 SOL) ---")
	t.Logf("Actual In (lamports): %s", quoteY.InAmount.String())
	t.Logf("Expected Out (USDC units): %s", quoteY.OutAmount.String())
	t.Logf("Min Out After Slippage (USDC units): %s", quoteY.MinOutAmount.String())
	t.Logf("Fee Paid (lamports): %s", quoteY.Fee.String())
	t.Logf("Price Impact: %f%%", quoteY.PriceImpact)
	t.Logf("Starting Active Bin ID: %d, Final Active Bin ID: %d", pair.ActiveID, quoteY.LastFilledActiveBinID)

	// 4. Simulating Swap Exact-In: 100 USDC -> SOL (swapForY = false)
	// Swap exact in: 100 USDC -> SOL
	inAmountUSDC := new(big.Int).SetUint64(100_000_000) // 100 USDC (6 decimals)
	quoteX, err := onchain.ComputeSwapQuote(pair, binArrays, inAmountUSDC, false, false, 50, 0, nil, nil)
	if err != nil {
		t.Fatalf("Failed to compute swap quote X: %v", err)
	}

	t.Logf("--- USDC -> SOL Swap Quote (100 USDC) ---")
	t.Logf("Actual In (USDC units): %s", quoteX.InAmount.String())
	t.Logf("Expected Out (SOL lamports): %s", quoteX.OutAmount.String())
	t.Logf("Min Out After Slippage (SOL lamports): %s", quoteX.MinOutAmount.String())
	t.Logf("Fee Paid (USDC units): %s", quoteX.Fee.String())
	t.Logf("Price Impact: %f%%", quoteX.PriceImpact)
	t.Logf("Starting Active Bin ID: %d, Final Active Bin ID: %d", pair.ActiveID, quoteX.LastFilledActiveBinID)
}

func BenchmarkComputeSwapQuote(b *testing.B) {
	c := onchain.NewClient(rpc.MainNetBeta_RPC)
	ctx := context.Background()
	poolAddr := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")

	pair, err := c.GetLBPair(ctx, poolAddr)
	if err != nil {
		b.Fatalf("Failed to fetch LBPair: %v", err)
	}

	activeIdx := onchain.BinIDToBinArrayIndex(pair.ActiveID)
	indices := []int64{activeIdx - 1, activeIdx, activeIdx + 1}
	var binArrays []*onchain.BinArray

	for _, idx := range indices {
		ba, err := c.GetBinArrayByIndex(ctx, poolAddr, idx)
		if err == nil {
			binArrays = append(binArrays, ba)
		}
	}

	inAmount := new(big.Int).SetUint64(1_000_000_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = onchain.ComputeSwapQuote(pair, binArrays, inAmount, true, false, 50, 0, nil, nil)
	}
}
