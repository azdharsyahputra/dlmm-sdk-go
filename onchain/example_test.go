package onchain_test

import (
	"fmt"
	"log"
	"math/big"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/gagliardetto/solana-go"
)

func ExampleDecodeLBPair() {
	var pair onchain.LBPair
	err := onchain.DecodeLBPair([]byte{1, 2, 3}, &pair)
	if err != nil {
		// handle err
	}

	fmt.Printf("Fetched Pool! Active Bin ID: %d\n", pair.ActiveID)
}

func ExampleComputeSwapQuote() {
	pair := &onchain.LBPair{
		ActiveID: -100,
		BinStep:  10,
	}
	binArrays := []*onchain.BinArray{}
	inAmount := new(big.Int).SetUint64(1_000_000)
	
	quote, err := onchain.ComputeSwapQuote(
		pair, binArrays, inAmount, true, false, 50, 0, nil, nil,
	)
	if err != nil {
		log.Printf("Failed to compute quote: %v", err)
		return
	}
	
	fmt.Printf("Expected Out Amount: %s\n", quote.OutAmount.String())
}

func ExampleNewSwapInstruction() {
	poolAddress := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")
	userWallet := solana.MustPublicKeyFromBase58("2wT8Yq49kHgDzFolAW5dK3a7R228G7k4A2m8j1x6zZ1x")
	
	params := onchain.SwapParams{
		LBPair:            poolAddress,
		ReserveX:          solana.NewWallet().PublicKey(),
		ReserveY:          solana.NewWallet().PublicKey(),
		UserTokenIn:       userWallet,                     
		UserTokenOut:      userWallet,                     
		TokenXMint:        solana.NewWallet().PublicKey(), 
		TokenYMint:        solana.NewWallet().PublicKey(), 
		Oracle:            solana.NewWallet().PublicKey(), 
		User:              userWallet,
		TokenXProgram:     solana.TokenProgramID,
		TokenYProgram:     solana.TokenProgramID,
		BinArrays:         []solana.PublicKey{solana.NewWallet().PublicKey()},
		AmountIn:          1_000_000,
		MinAmountOut:      990_000,
	}

	inst, err := onchain.NewSwapInstruction(params)
	if err != nil {
		log.Fatalf("Failed to build swap instruction: %v", err)
	}

	fmt.Printf("Instruction Program ID: %s\n", inst.ProgramID().String())
}

func ExampleBinIDToTokenPrice() {
	binID := int32(-6368)
	binStep := uint16(4)
	decimalsX := int32(9)
	decimalsY := int32(6)

	price := onchain.BinIDToTokenPrice(binID, binStep, decimalsX, decimalsY)
	fmt.Printf("Price of SOL at Bin %d is approx %.2f USDC\n", binID, price)
}

func ExampleNewClaimFeeInstruction() {
	poolAddress := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")
	userWallet := solana.MustPublicKeyFromBase58("2wT8Yq49kHgDzFolAW5dK3a7R228G7k4A2m8j1x6zZ1x")
	position := solana.MustPublicKeyFromBase58("Pos1111111111111111111111111111111111111111")
	
	params := onchain.ClaimFeeParams{
		LBPair:        poolAddress,
		Position:      position,
		Sender:        userWallet,
		ReserveX:      solana.NewWallet().PublicKey(),
		ReserveY:      solana.NewWallet().PublicKey(),
		UserTokenX:    solana.NewWallet().PublicKey(),
		UserTokenY:    solana.NewWallet().PublicKey(),
		TokenXMint:    solana.NewWallet().PublicKey(),
		TokenYMint:    solana.NewWallet().PublicKey(),
		TokenProgramX: solana.TokenProgramID,
		TokenProgramY: solana.TokenProgramID,
	}

	inst, err := onchain.NewClaimFeeInstruction(params)
	if err != nil {
		log.Fatalf("Failed to build claim fee instruction: %v", err)
	}

	fmt.Printf("Instruction Program ID: %s\n", inst.ProgramID().String())
}

func ExampleNewClaimRewardInstruction() {
	poolAddress := solana.MustPublicKeyFromBase58("5rCf1DM8LjKTw4YqhnoLcngyZYeNnQqztScTogYHAS6")
	userWallet := solana.MustPublicKeyFromBase58("2wT8Yq49kHgDzFolAW5dK3a7R228G7k4A2m8j1x6zZ1x")
	position := solana.MustPublicKeyFromBase58("Pos1111111111111111111111111111111111111111")
	
	params := onchain.ClaimRewardParams{
		LBPair:           poolAddress,
		Position:         position,
		Sender:           userWallet,
		RewardVault:      solana.NewWallet().PublicKey(),
		RewardMint:       solana.NewWallet().PublicKey(),
		UserTokenAccount: solana.NewWallet().PublicKey(),
		TokenProgram:     solana.TokenProgramID,
		RewardIndex:      0,
		MinBinID:         -10,
		MaxBinID:         10,
	}

	inst, err := onchain.NewClaimRewardInstruction(params)
	if err != nil {
		log.Fatalf("Failed to build claim reward instruction: %v", err)
	}

	fmt.Printf("Instruction Program ID: %s\n", inst.ProgramID().String())
}
