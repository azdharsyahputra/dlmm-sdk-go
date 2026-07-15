package client_test

import (
	"context"
	"fmt"
	"log"

	"github.com/azdharsyahputra/dlmm-sdk-go/client"
)

func ExampleNewClient() {
	// Initialize a new DLMM REST API client targeting mainnet (default)
	apiClient := client.NewClient("")

	// Or specify a custom devnet endpoint
	devClient := client.NewClient("https://dlmm.dev.metdev.io")
	_ = devClient
	
	fmt.Printf("Client initialized: %T\n", apiClient)
	// Output: Client initialized: *client.Client
}

func ExampleClient_GetPools() {
	apiClient := client.NewClient("")
	ctx := context.Background()

	page := 1
	pageSize := 1
	poolsRes, err := apiClient.GetPools(ctx, &client.GetPoolsParams{
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(poolsRes.Data) > 0 {
		fmt.Printf("Successfully fetched pool: %s\n", poolsRes.Data[0].Name)
	}
}

func ExampleClient_GetPool() {
	apiClient := client.NewClient("")
	ctx := context.Background()
	
	poolAddress := "5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6"
	pool, err := apiClient.GetPool(ctx, poolAddress)
	if err != nil {
		log.Fatalf("Failed to fetch pool details: %v", err)
	}

	fmt.Printf("Pool Name: %s\n", pool.Name)
}

func ExampleClient_GetPortfolio() {
	apiClient := client.NewClient("")
	ctx := context.Background()

	walletAddress := "2wT8Yq49kHgDzFolAW5dK3a7R228G7k4A2m8j1x6zZ1x"

	portfolio, err := apiClient.GetPortfolio(ctx, walletAddress, &client.GetPortfolioParams{})
	if err != nil {
		log.Fatalf("Failed to fetch portfolio: %v", err)
	}

	fmt.Printf("Fetched portfolio for %s, active pools count: %d\n", walletAddress, len(portfolio.Pools))
}

func ExampleClient_GetOHLCV() {
	apiClient := client.NewClient("")
	ctx := context.Background()

	poolAddress := "5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6"
	tf := "15m"
	
	res, err := apiClient.GetOHLCV(ctx, poolAddress, &client.GetOHLCVParams{
		Timeframe: &tf,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(res.Data) > 0 {
		fmt.Printf("First OHLCV Close Price: %f\n", res.Data[0].Close)
	}
}

func ExampleClient_GetWalletPoolTotalClaims() {
	apiClient := client.NewClient("")
	ctx := context.Background()

	walletAddress := "2wT8Yq49kHgDzFolAW5dK3a7R228G7k4A2m8j1x6zZ1x"
	poolAddress := "5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6"

	claims, err := apiClient.GetWalletPoolTotalClaims(ctx, walletAddress, poolAddress)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Total Claimed Fee (Token X): %s\n", claims.ClaimedFeeX)
}

func ExampleClient_GetOpenPortfolio() {
	apiClient := client.NewClient("")
	ctx := context.Background()

	res, err := apiClient.GetOpenPortfolio(ctx, &client.GetOpenPortfolioParams{})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Found %d open limit orders.\n", len(res.Pools))
}

func ExampleClient_GetPortfolioTotal() {
	apiClient := client.NewClient("")
	ctx := context.Background()

	walletAddress := "2wT8Yq49kHgDzFolAW5dK3a7R228G7k4A2m8j1x6zZ1x"
	
	res, err := apiClient.GetPortfolioTotal(ctx, walletAddress)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Total value of portfolio in USD: %s\n", res.TotalPnlUsd)
}
