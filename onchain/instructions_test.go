package onchain_test

import (
	"testing"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
)

func TestSwapInstructionBuilder(t *testing.T) {
	lbPair := solana.NewWallet().PublicKey()
	reserveX := solana.NewWallet().PublicKey()
	reserveY := solana.NewWallet().PublicKey()
	userTokenIn := solana.NewWallet().PublicKey()
	userTokenOut := solana.NewWallet().PublicKey()
	tokenXMint := solana.NewWallet().PublicKey()
	tokenYMint := solana.NewWallet().PublicKey()
	oracle := solana.NewWallet().PublicKey()
	user := solana.NewWallet().PublicKey()

	tokenProgram := solana.TokenProgramID

	ba1 := solana.NewWallet().PublicKey()
	ba2 := solana.NewWallet().PublicKey()

	params := onchain.SwapParams{
		LbPair:        lbPair,
		ReserveX:      reserveX,
		ReserveY:      reserveY,
		UserTokenIn:   userTokenIn,
		UserTokenOut:  userTokenOut,
		TokenXMint:    tokenXMint,
		TokenYMint:    tokenYMint,
		Oracle:        oracle,
		User:          user,
		TokenXProgram: tokenProgram,
		TokenYProgram: tokenProgram,
		BinArrays:     []solana.PublicKey{ba1, ba2},
		AmountIn:      1000000,
		MinAmountOut:  990000,
	}

	ix, err := onchain.NewSwapInstruction(params)
	if err != nil {
		t.Fatalf("Failed to build swap instruction: %v", err)
	}

	if !ix.ProgramID().Equals(onchain.ProgramID) {
		t.Errorf("Expected program ID %s, got %s", onchain.ProgramID.String(), ix.ProgramID().String())
	}

	expectedAccountsCount := 18
	if len(ix.Accounts()) != expectedAccountsCount {
		t.Errorf("Expected %d accounts, got %d", expectedAccountsCount, len(ix.Accounts()))
	}

	data, err := ix.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}
	expectedDataLen := 28
	if len(data) != expectedDataLen {
		t.Errorf("Expected data length %d, got %d", expectedDataLen, len(data))
	}

	t.Log("Successfully verified swap instruction builder!")
}

func TestClaimFeeInstructionBuilder(t *testing.T) {
	params := onchain.ClaimFeeParams{
		LbPair:        solana.NewWallet().PublicKey(),
		Position:      solana.NewWallet().PublicKey(),
		Sender:        solana.NewWallet().PublicKey(),
		ReserveX:      solana.NewWallet().PublicKey(),
		ReserveY:      solana.NewWallet().PublicKey(),
		UserTokenX:    solana.NewWallet().PublicKey(),
		UserTokenY:    solana.NewWallet().PublicKey(),
		TokenXMint:    solana.NewWallet().PublicKey(),
		TokenYMint:    solana.NewWallet().PublicKey(),
		TokenProgramX: solana.TokenProgramID,
		TokenProgramY: solana.TokenProgramID,
		MinBinId:      -100,
		MaxBinId:      100,
	}

	ix, err := onchain.NewClaimFeeInstruction(params)
	if err != nil {
		t.Fatalf("Failed to build claim fee instruction: %v", err)
	}

	if !ix.ProgramID().Equals(onchain.ProgramID) {
		t.Errorf("Expected program ID %s, got %s", onchain.ProgramID.String(), ix.ProgramID().String())
	}

	// 14 static accounts
	expectedAccountsCount := 14
	if len(ix.Accounts()) != expectedAccountsCount {
		t.Errorf("Expected %d accounts, got %d", expectedAccountsCount, len(ix.Accounts()))
	}

	// data length: 8 (disc) + 4 (min bin) + 4 (max bin) + 4 (slices len vec) = 20 bytes
	data, err := ix.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}
	expectedDataLen := 20
	if len(data) != expectedDataLen {
		t.Errorf("Expected data length %d, got %d", expectedDataLen, len(data))
	}

	t.Log("Successfully verified claim fee instruction builder!")
}

func TestClaimRewardInstructionBuilder(t *testing.T) {
	params := onchain.ClaimRewardParams{
		LbPair:           solana.NewWallet().PublicKey(),
		Position:         solana.NewWallet().PublicKey(),
		Sender:           solana.NewWallet().PublicKey(),
		RewardVault:      solana.NewWallet().PublicKey(),
		RewardMint:       solana.NewWallet().PublicKey(),
		UserTokenAccount: solana.NewWallet().PublicKey(),
		TokenProgram:     solana.TokenProgramID,
		RewardIndex:      0,
		MinBinId:         -100,
		MaxBinId:         100,
	}

	ix, err := onchain.NewClaimRewardInstruction(params)
	if err != nil {
		t.Fatalf("Failed to build claim reward instruction: %v", err)
	}

	if !ix.ProgramID().Equals(onchain.ProgramID) {
		t.Errorf("Expected program ID %s, got %s", onchain.ProgramID.String(), ix.ProgramID().String())
	}

	// 10 static accounts
	expectedAccountsCount := 10
	if len(ix.Accounts()) != expectedAccountsCount {
		t.Errorf("Expected %d accounts, got %d", expectedAccountsCount, len(ix.Accounts()))
	}

	// data length: 8 (disc) + 8 (reward idx) + 4 (min bin) + 4 (max bin) + 4 (slices len vec) = 28 bytes
	data, err := ix.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}
	expectedDataLen := 28
	if len(data) != expectedDataLen {
		t.Errorf("Expected data length %d, got %d", expectedDataLen, len(data))
	}

	t.Log("Successfully verified claim reward instruction builder!")
}
