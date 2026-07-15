package client_test

import (
	"context"
	"fmt"
	"log"

	"github.com/azdharsyahputra/dlmm-sdk-go/client"
)

func ExampleNewClient() {
	// Initialize a new DLMM REST API client
	apiClient := client.NewClient("")

	// Example: Fetch details for the SOL-USDC pool
	ctx := context.Background()
	poolAddress := "5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6"
	
	poolDetails, err := apiClient.GetPoolDetails(ctx, poolAddress)
	if err != nil {
		log.Fatalf("Failed to fetch pool details: %v", err)
	}

	fmt.Printf("Pool Name: %s\n", poolDetails.Name)
	fmt.Printf("Current Price: %f\n", poolDetails.CurrentPrice)
}

func ExampleClient_GetPairs() {
	apiClient := client.NewClient("")
	ctx := context.Background()
	
	// Fetch pairs with pagination (limit 5)
	pairs, err := apiClient.GetPairs(ctx, 5, 0)
	if err != nil {
		log.Fatalf("Failed to fetch pairs: %v", err)
	}

	if len(pairs) > 0 {
		fmt.Printf("Successfully fetched %d pools. First pool: %s\n", len(pairs), pairs[0].Name)
	}
}
