package main

import (
	"context"
	"fmt"
	"log"

	"github.com/azdharsyahputra/dlmm-sdk-go/client"
	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	ctx := context.Background()

	fmt.Println("1. Connecting to Devnet API (https://dlmm.dev.metdev.io)...")
	devnetAPI := client.NewClient("https://dlmm.dev.metdev.io")

	page := 1
	pageSize := 1
	poolsRes, err := devnetAPI.GetPools(ctx, &client.GetPoolsParams{
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatalf("Failed to fetch pairs: %v", err)
	}
	if len(poolsRes.Data) == 0 {
		log.Fatalf("No pools found on devnet")
	}

	targetPool := poolsRes.Data[0]
	fmt.Printf("   Found Devnet Pool: %s (Address: %s)\n", targetPool.Name, targetPool.Address)

	fmt.Println("\n2. Connecting to Solana Mainnet RPC (because dev.metdev.io indexes mainnet)...")
	rpcClient := rpc.New(rpc.MainNetBeta_RPC)
	poolPubkey := solana.MustPublicKeyFromBase58(targetPool.Address)

	fmt.Println("   Fetching On-Chain LbPair Data...")
	accountInfo, err := rpcClient.GetAccountInfo(ctx, poolPubkey)
	if err != nil {
		log.Fatalf("Failed to fetch LbPair account: %v", err)
	}

	fmt.Println("\n3. Deserializing raw bytes into Go struct...")
	var pair onchain.LbPair
	err = onchain.DeserializeAccount(accountInfo.Value.Data.GetBinary(), &pair)
	if err != nil {
		log.Fatalf("Failed to decode LbPair: %v", err)
	}

	fmt.Println("   Decode Success!")
	fmt.Printf("   - Pool Active Bin ID : %d\n", pair.ActiveId)
	fmt.Printf("   - Bin Step           : %d\n", pair.BinStep)
	fmt.Printf("   - Reserve X          : %s\n", pair.ReserveX.String())
	fmt.Printf("   - Reserve Y          : %s\n", pair.ReserveY.String())
	
	fmt.Println("\n✅ ALL SYSTEMS GO! SDK is perfectly interacting with Meteora Devnet.")
}
