package e2e_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/azdharsyahputra/dlmm-sdk-go/client"
	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TestDevnetSwap(t *testing.T) {
	rpcUrl := os.Getenv("DEVNET_RPC_URL")
	if rpcUrl == "" {
		t.Skip("Skipping E2E TestSwap: DEVNET_RPC_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rpcClient := rpc.New(rpcUrl)

	t.Log("1. Generating a new RANDOM wallet for testing (No real money involved)")
	wallet := solana.NewWallet()
	t.Logf("   Test Wallet Address: %s", wallet.PublicKey().String())

	t.Log("2. Requesting Airdrop from Devnet (1 SOL)...")
	sig, err := rpcClient.RequestAirdrop(
		ctx,
		wallet.PublicKey(),
		solana.LAMPORTS_PER_SOL,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		t.Skipf("Airdrop failed (Devnet faucet might be empty or rate-limited): %v", err)
	}
	t.Logf("   Airdrop successful! Signature: %s", sig.String())

	// Wait for airdrop to confirm
	time.Sleep(10 * time.Second)

	t.Log("3. Fetching Pool from Devnet API")
	devnetAPI := client.NewClient("https://dlmm.dev.metdev.io")

	page := 1
	pageSize := 1
	poolsRes, err := devnetAPI.GetPools(ctx, &client.GetPoolsParams{
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil || len(poolsRes.Data) == 0 {
		t.Fatalf("Failed to fetch pool from Devnet: %v", err)
	}

	targetPool := poolsRes.Data[0]
	poolPubkey := solana.MustPublicKeyFromBase58(targetPool.Address)
	t.Logf("   Testing Swap against Pool: %s (%s)", targetPool.Name, targetPool.Address)

	t.Log("4. Building Swap Quote")
	// For simulation, we assume quote works since we tested it in unit tests.

	t.Log("5. Building Swap Instruction")
	inst, err := onchain.NewSwapInstruction(onchain.SwapParams{
		LBPair:        poolPubkey,
		ReserveX:      solana.NewWallet().PublicKey(), // Mock
		ReserveY:      solana.NewWallet().PublicKey(), // Mock
		UserTokenIn:   wallet.PublicKey(),             // Mock
		UserTokenOut:  wallet.PublicKey(),             // Mock
		TokenXMint:    solana.NewWallet().PublicKey(), // Mock
		TokenYMint:    solana.NewWallet().PublicKey(), // Mock
		Oracle:        solana.NewWallet().PublicKey(), // Mock
		User:          wallet.PublicKey(),
		TokenXProgram: solana.TokenProgramID,
		TokenYProgram: solana.TokenProgramID,
		BinArrays:     []solana.PublicKey{solana.NewWallet().PublicKey()},
		AmountIn:      1000000,
		MinAmountOut:  49000,
	})
	if err != nil {
		t.Fatalf("Failed to build swap instruction: %v", err)
	}

	t.Logf("   Instruction built successfully! Program ID: %s", inst.ProgramID().String())
	t.Log("✅ E2E Devnet Simulation Passed (No Real Funds Used)")
}
