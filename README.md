# Meteora DLMM Go SDK (`github.com/azdharsyahputra/dlmm-sdk-go`)

A complete Go library for interacting with **Meteora's Concentrated Liquidity Market Maker (DLMM)** protocol. 

This SDK provides two components:
1. **REST API Client** (`client` package): Queries historical data, statistics, active/closed portfolios, position PnL, daily volume/fees, and limit orders from Meteora's analytics endpoint.
2. **On-Chain Client** (`onchain` package): Interfaces directly with the Solana blockchain to fetch pool accounts (`LbPair`, `BinArray`, `PositionV2`), derive PDA accounts, perform off-chain swap quoting/simulations, and construct instruction transactions.

---

## Installation

```bash
go get github.com/azdharsyahputra/dlmm-sdk-go
```

Make sure your Go project initializes modules and fetches dependencies:
```bash
go mod tidy
```

---

## 1. REST API Client Usage

The REST API client communicates with `https://dlmm.datapi.meteora.ag` to retrieve statistical and analytical data.

### Example: Fetching Pools and Protocol Overview

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/azdharsyahputra/dlmm-sdk-go/client"
)

func main() {
	// Initialize the client (defaults to Meteora DLMM production API)
	c := client.NewClient("")
	ctx := context.Background()

	// Fetch top 5 pools
	pageSize := 5
	pools, err := c.GetPools(ctx, &client.GetPoolsParams{
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatalf("Failed to fetch pools: %v", err)
	}

	for _, pool := range pools.Data {
		fmt.Printf("Pool: %s | Address: %s | TVL: $%.2f | APR: %.2f%%\n",
			pool.Name, pool.Address, pool.Tvl, pool.Apr)
	}

	// Fetch protocol overview stats
	overview, err := c.GetProtocolOverview(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch overview: %v", err)
	}
	fmt.Printf("Meteora DLMM Total 24h Volume: $%.2f\n", overview.Volume24h)
}
```

---

## 2. On-Chain Client Usage

The on-chain client uses `github.com/gagliardetto/solana-go` to connect to a Solana RPC node and read state or build transactions.

### Example: Fetching Pool State and Simulating a Swap Quote

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	// Connect to Solana Mainnet RPC
	c := onchain.NewClient(rpc.MainNetBeta_RPC)
	ctx := context.Background()

	// SOL-USDC DLMM pool address
	poolAddr := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")

	// 1. Fetch live pool state (LbPair)
	pair, err := c.GetLbPair(ctx, poolAddr)
	if err != nil {
		log.Fatalf("Failed to fetch pool: %v", err)
	}
	fmt.Printf("Active Bin ID: %d | Bin Step: %d BPS\n", pair.ActiveId, pair.BinStep)

	// Calculate current SOL price (SOL = 9 decimals, USDC = 6 decimals)
	price := onchain.BinIdToTokenPrice(pair.ActiveId, pair.BinStep, 9, 6)
	fmt.Printf("Current price: %.4f USDC per SOL\n", price)

	// 2. Fetch the active BinArray containing the current price
	activeIdx := onchain.BinIdToBinArrayIndex(pair.ActiveId)
	binArray, err := c.GetBinArrayByIndex(ctx, poolAddr, activeIdx)
	if err != nil {
		log.Fatalf("Failed to fetch active bin array: %v", err)
	}

	// 3. Simulating swap: 1 SOL -> USDC (swapForY = true)
	solInAmount := big.NewInt(1_000_000_000) // 1 SOL in lamports
	quote, err := onchain.ComputeSwapQuote(pair, []onchain.BinArray{*binArray}, solInAmount, true, false, 50, 0)
	if err != nil {
		log.Fatalf("Swap quote failed: %v", err)
	}

	fmt.Printf("--- Swap Quote: 1 SOL to USDC ---\n")
	fmt.Printf("Expected Out: %s USDC units (~$%.2f)\n", quote.OutAmount.String(), float64(quote.OutAmount.Int64())/1e6)
	fmt.Printf("Fee Paid: %s lamports\n", quote.Fee.String())
	fmt.Printf("Price Impact: %.6f%%\n", quote.PriceImpact)
}
```

### Example: Building a Swap Instruction

```go
package main

import (
	"fmt"
	"log"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
)

func main() {
	// Construct the instructions params
	params := onchain.SwapParams{
		LbPair:        solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6"),
		ReserveX:      solana.MustPublicKeyFromBase58("..."), // reserve vaults derived from LbPair
		ReserveY:      solana.MustPublicKeyFromBase58("..."),
		UserTokenIn:   solana.MustPublicKeyFromBase58("..."), // user token accounts
		UserTokenOut:  solana.MustPublicKeyFromBase58("..."),
		TokenXMint:    solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"),
		TokenYMint:    solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
		Oracle:        solana.MustPublicKeyFromBase58("..."), // derived from LbPair
		User:          solana.MustPublicKeyFromBase58("..."), // signer wallet
		TokenXProgram: solana.TokenProgramID,
		TokenYProgram: solana.TokenProgramID,
		BinArrays:     []solana.PublicKey{
			// bin arrays to be traversed derived from ComputeSwapQuote's final active bin ID
		},
		AmountIn:     1000000000,
		MinAmountOut: 990000000, // custom slippage applied to Quote.OutAmount
	}

	swapIx, err := onchain.NewSwapInstruction(params)
	if err != nil {
		log.Fatalf("Failed to build swap instruction: %v", err)
	}

	fmt.Println("Created swap instruction successfully. Program ID:", swapIx.ProgramID())
}
```

### High-Precision Decimal Math

To avoid float rounding issues in production, use the `decimal.Decimal` variants for bin-to-price conversions:

```go
import (
	"fmt"
	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
)

func main() {
	// High-precision token price (Y per X)
	priceDec := onchain.BinIdToTokenPriceDecimal(-6388, 4, 9, 6)
	fmt.Println("SOL-USDC Price (Decimal):", priceDec.String()) // "77.7163906299..."

	// High-precision price-to-bin conversion
	binId := onchain.TokenPriceToBinIdDecimal(priceDec, 4, 9, 6)
	fmt.Println("Derived Bin ID:", binId) // -6388
}
```

---

## Running Tests

To run the integration and unit tests (uses live mainnet calls to verify correctness):

```bash
# Test REST API Client
go test -v ./client

# Test On-Chain Client (Layouts, Quoting, Instructions)
go test -v ./onchain
```
