// Package client provides a REST API client for Meteora's DLMM architecture.
//
// It allows users to interact with Meteora's off-chain APIs to fetch pool data,
// historical analytics (OHLCV, Volume), portfolio positions, and daily protocol metrics.
//
// The client supports optional configuration via the Functional Options pattern, 
// such as injecting a custom *http.Client for timeouts or proxying.
//
// # Quick Start
//
//	apiClient := client.NewClient("") // Uses default mainnet endpoint
//	ctx := context.Background()
//
//	pool, err := apiClient.GetPool(ctx, "5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Pool TVL:", pool.Tvl)
package client
