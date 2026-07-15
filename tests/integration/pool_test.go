package integration_test

import (
	"context"
	"testing"

	"github.com/azdharsyahputra/dlmm-sdk-go/client"
)

func TestLiveAPI(t *testing.T) {
	// Initialize live client
	c := client.NewClient("")

	ctx := context.Background()

	// 1. Test GetPools
	pageSize := 5
	poolsResp, err := c.GetPools(ctx, &client.GetPoolsParams{
		PageSize: &pageSize,
	})
	if err != nil {
		t.Fatalf("Failed to get pools: %v", err)
	}

	if len(poolsResp.Data) == 0 {
		t.Fatal("Expected at least one pool, got 0")
	}

	t.Logf("Successfully fetched %d pools. First pool: %s (%s)", len(poolsResp.Data), poolsResp.Data[0].Name, poolsResp.Data[0].Address)

	// 2. Test GetPool for the first pool
	poolAddr := poolsResp.Data[0].Address
	pool, err := c.GetPool(ctx, poolAddr)
	if err != nil {
		t.Fatalf("Failed to get pool details for %s: %v", poolAddr, err)
	}

	if pool.Address != poolAddr {
		t.Errorf("Expected pool address %s, got %s", poolAddr, pool.Address)
	}

	t.Logf("Successfully fetched pool details. TVL: %f, Current Price: %f", pool.Tvl, pool.CurrentPrice)

	// 3. Test GetProtocolOverview
	overview, err := c.GetProtocolOverview(ctx)
	if err != nil {
		t.Fatalf("Failed to get protocol overview: %v", err)
	}

	t.Logf("Protocol Overview: TVL: %f, 24h Volume: %f, 24h Fees: %f", overview.Tvl, overview.Volume24h, overview.Fees24h)
}
