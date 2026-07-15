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

	// 3. Get pool details
	ctx := context.Background()
	pool, err := apiClient.GetPool(ctx, "5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")
	if err != nil {
		log.Fatalf("Failed to fetch pool details: %v", err)
	}

	fmt.Printf("Pool Name: %s\n", pool.Name)
	fmt.Printf("Current Price: %f\n", pool.CurrentPrice)
}

func ExampleClient_GetPools() {
	clientAPI := client.NewClient("") // Use default mainnet URL
	ctx := context.Background()

	page := 1
	pageSize := 1
	poolsRes, err := clientAPI.GetPools(ctx, &client.GetPoolsParams{
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, p := range poolsRes.Data {
		fmt.Printf("Pool: %s, Address: %s\n", p.Name, p.Address)
	}
}
