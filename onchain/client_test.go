package onchain_test

import (
	"context"
	"testing"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TestLiveOnChainLBKeyPair(t *testing.T) {
	// Initialize mainnet client
	// We can use the public Solana mainnet RPC or Meteora's
	c := onchain.NewClient(rpc.MainNetBeta_RPC)

	ctx := context.Background()

	// SOL-USDC pool address on mainnet
	poolAddr := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")

	// 1. Fetch LbPair
	pair, err := c.GetLbPair(ctx, poolAddr)
	if err != nil {
		t.Fatalf("Failed to fetch LbPair: %v", err)
	}

	t.Logf("LbPair Token X: %s", pair.TokenXMint.String())
	t.Logf("LbPair Token Y: %s", pair.TokenYMint.String())
	t.Logf("LbPair Active Bin ID: %d", pair.ActiveId)
	t.Logf("LbPair Bin Step: %d BPS", pair.BinStep)

	// Verify token mints are not empty public keys
	if pair.TokenXMint.IsZero() || pair.TokenYMint.IsZero() {
		t.Error("Expected token mint public keys to be non-zero")
	}

	// 2. Fetch active BinArray
	binArrayIndex := onchain.BinIdToBinArrayIndex(pair.ActiveId)
	t.Logf("Active Bin ID %d is in BinArray index %d", pair.ActiveId, binArrayIndex)

	binArray, err := c.GetBinArrayByIndex(ctx, poolAddr, binArrayIndex)
	if err != nil {
		t.Fatalf("Failed to fetch BinArray: %v", err)
	}

	if binArray.LbPair != poolAddr {
		t.Errorf("BinArray LbPair address mismatch: expected %s, got %s", poolAddr.String(), binArray.LbPair.String())
	}

	t.Logf("BinArray Index: %d", binArray.Index)

	// 3. Extract active Bin
	activeBin, err := onchain.GetBinFromBinArrayHelper(pair.ActiveId, binArray)
	if err != nil {
		t.Fatalf("Failed to get active bin from BinArray: %v", err)
	}

	t.Logf("Active Bin X Amount: %d, Y Amount: %d", activeBin.AmountX, activeBin.AmountY)

	// 4. Calculate prices
	lamportPrice := onchain.BinIdToPrice(pair.ActiveId, pair.BinStep)
	tokenPrice := onchain.BinIdToTokenPrice(pair.ActiveId, pair.BinStep, 9, 6) // SOL is 9 decimals, USDC is 6 decimals
	t.Logf("Calculated Lamport Price: %f", lamportPrice)
	t.Logf("Calculated SOL-USDC Token Price: %f (USDC per SOL)", tokenPrice)

	if tokenPrice <= 0 {
		t.Error("Expected token price to be greater than 0")
	}
}
