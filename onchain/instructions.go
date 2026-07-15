package onchain

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// DeriveEventAuthority derives the event authority PDA for DLMM.
func DeriveEventAuthority() (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("__event_authority"),
	}
	return solana.FindProgramAddress(seeds, ProgramID)
}

// SwapParams defines parameters to build a Swap instruction.
type SwapParams struct {
	LbPair                   solana.PublicKey
	ReserveX                 solana.PublicKey
	ReserveY                 solana.PublicKey
	UserTokenIn              solana.PublicKey
	UserTokenOut             solana.PublicKey
	TokenXMint               solana.PublicKey
	TokenYMint               solana.PublicKey
	Oracle                   solana.PublicKey
	BinArrayBitmapExtension  *solana.PublicKey
	HostFeeIn                *solana.PublicKey
	User                     solana.PublicKey
	TokenXProgram            solana.PublicKey
	TokenYProgram            solana.PublicKey
	BinArrays                []solana.PublicKey
	AmountIn                 uint64
	MinAmountOut             uint64
}

// NewSwapInstruction builds a swap2 instruction matching the DLMM program layout.
func NewSwapInstruction(params SwapParams) (*solana.GenericInstruction, error) {
	eventAuthority, _, err := DeriveEventAuthority()
	if err != nil {
		return nil, fmt.Errorf("failed to derive event authority PDA: %w", err)
	}

	memoProgram := solana.MustPublicKeyFromBase58("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")

	bitmapExt := ProgramID
	if params.BinArrayBitmapExtension != nil {
		bitmapExt = *params.BinArrayBitmapExtension
	}

	hostFeeIn := ProgramID
	if params.HostFeeIn != nil {
		hostFeeIn = *params.HostFeeIn
	}

	// Define accounts list for swap2
	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(params.LbPair, true, false), // lb_pair
		solana.NewAccountMeta(bitmapExt, true, false),      // bin_array_bitmap_extension
		solana.NewAccountMeta(params.ReserveX, true, false), // reserve_x
		solana.NewAccountMeta(params.ReserveY, true, false), // reserve_y
		solana.NewAccountMeta(params.UserTokenIn, true, false), // user_token_in
		solana.NewAccountMeta(params.UserTokenOut, true, false), // user_token_out
		solana.NewAccountMeta(params.TokenXMint, false, false), // token_x_mint
		solana.NewAccountMeta(params.TokenYMint, false, false), // token_y_mint
		solana.NewAccountMeta(params.Oracle, true, false), // oracle
		solana.NewAccountMeta(hostFeeIn, true, false),      // host_fee_in
		solana.NewAccountMeta(params.User, false, true),   // user
		solana.NewAccountMeta(params.TokenXProgram, false, false), // token_x_program
		solana.NewAccountMeta(params.TokenYProgram, false, false), // token_y_program
		solana.NewAccountMeta(memoProgram, false, false),          // memo_program
		solana.NewAccountMeta(eventAuthority, false, false),       // event_authority
		solana.NewAccountMeta(ProgramID, false, false),            // program (DLMM itself)
	}

	// Append bin array accounts as remaining accounts
	for _, ba := range params.BinArrays {
		accounts = append(accounts, solana.NewAccountMeta(ba, true, false))
	}

	// Serialize data
	buf := new(bytes.Buffer)
	// 1. Discriminator for swap2
	discriminator := []byte{65, 75, 63, 76, 235, 91, 91, 136}
	buf.Write(discriminator)

	// 2. amount_in (u64)
	if err := binary.Write(buf, binary.LittleEndian, params.AmountIn); err != nil {
		return nil, fmt.Errorf("failed to serialize amount_in: %w", err)
	}

	// 3. min_amount_out (u64)
	if err := binary.Write(buf, binary.LittleEndian, params.MinAmountOut); err != nil {
		return nil, fmt.Errorf("failed to serialize min_amount_out: %w", err)
	}

	// 4. remaining_accounts_info (empty slices vec = [0, 0, 0, 0] in Borsh)
	buf.Write([]byte{0, 0, 0, 0})

	return solana.NewInstruction(ProgramID, accounts, buf.Bytes()), nil
}

// ClaimFeeParams defines parameters to build a ClaimFee2 instruction.
type ClaimFeeParams struct {
	LbPair        solana.PublicKey
	Position      solana.PublicKey
	Sender        solana.PublicKey
	ReserveX      solana.PublicKey
	ReserveY      solana.PublicKey
	UserTokenX    solana.PublicKey
	UserTokenY    solana.PublicKey
	TokenXMint    solana.PublicKey
	TokenYMint    solana.PublicKey
	TokenProgramX solana.PublicKey
	TokenProgramY solana.PublicKey
	MinBinId      int32
	MaxBinId      int32
}

// NewClaimFeeInstruction builds a claim_fee2 instruction.
func NewClaimFeeInstruction(params ClaimFeeParams) (*solana.GenericInstruction, error) {
	eventAuthority, _, err := DeriveEventAuthority()
	if err != nil {
		return nil, fmt.Errorf("failed to derive event authority: %w", err)
	}

	memoProgram := solana.MustPublicKeyFromBase58("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")

	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(params.LbPair, true, false),
		solana.NewAccountMeta(params.Position, true, false),
		solana.NewAccountMeta(params.Sender, false, true), // signer
		solana.NewAccountMeta(params.ReserveX, true, false),
		solana.NewAccountMeta(params.ReserveY, true, false),
		solana.NewAccountMeta(params.UserTokenX, true, false),
		solana.NewAccountMeta(params.UserTokenY, true, false),
		solana.NewAccountMeta(params.TokenXMint, false, false),
		solana.NewAccountMeta(params.TokenYMint, false, false),
		solana.NewAccountMeta(params.TokenProgramX, false, false),
		solana.NewAccountMeta(params.TokenProgramY, false, false),
		solana.NewAccountMeta(memoProgram, false, false),
		solana.NewAccountMeta(eventAuthority, false, false),
		solana.NewAccountMeta(ProgramID, false, false),
	}

	buf := new(bytes.Buffer)
	// discriminator for claim_fee2: [112, 191, 101, 171, 28, 144, 127, 187]
	discriminator := []byte{112, 191, 101, 171, 28, 144, 127, 187}
	buf.Write(discriminator)

	binary.Write(buf, binary.LittleEndian, params.MinBinId)
	binary.Write(buf, binary.LittleEndian, params.MaxBinId)
	// empty slices vec
	buf.Write([]byte{0, 0, 0, 0})

	return solana.NewInstruction(ProgramID, accounts, buf.Bytes()), nil
}

// ClaimRewardParams defines parameters to build a ClaimReward2 instruction.
type ClaimRewardParams struct {
	LbPair           solana.PublicKey
	Position         solana.PublicKey
	Sender           solana.PublicKey
	RewardVault      solana.PublicKey
	RewardMint       solana.PublicKey
	UserTokenAccount solana.PublicKey
	TokenProgram     solana.PublicKey
	RewardIndex      uint64
	MinBinId         int32
	MaxBinId         int32
}

// NewClaimRewardInstruction builds a claim_reward2 instruction.
func NewClaimRewardInstruction(params ClaimRewardParams) (*solana.GenericInstruction, error) {
	eventAuthority, _, err := DeriveEventAuthority()
	if err != nil {
		return nil, fmt.Errorf("failed to derive event authority: %w", err)
	}

	memoProgram := solana.MustPublicKeyFromBase58("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")

	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(params.LbPair, true, false),
		solana.NewAccountMeta(params.Position, true, false),
		solana.NewAccountMeta(params.Sender, false, true), // signer
		solana.NewAccountMeta(params.RewardVault, true, false),
		solana.NewAccountMeta(params.RewardMint, false, false),
		solana.NewAccountMeta(params.UserTokenAccount, true, false),
		solana.NewAccountMeta(params.TokenProgram, false, false),
		solana.NewAccountMeta(memoProgram, false, false),
		solana.NewAccountMeta(eventAuthority, false, false),
		solana.NewAccountMeta(ProgramID, false, false),
	}

	buf := new(bytes.Buffer)
	// discriminator for claim_reward2: [190, 3, 127, 119, 178, 87, 157, 183]
	discriminator := []byte{190, 3, 127, 119, 178, 87, 157, 183}
	buf.Write(discriminator)

	binary.Write(buf, binary.LittleEndian, params.RewardIndex)
	binary.Write(buf, binary.LittleEndian, params.MinBinId)
	binary.Write(buf, binary.LittleEndian, params.MaxBinId)
	// empty slices vec
	buf.Write([]byte{0, 0, 0, 0})

	return solana.NewInstruction(ProgramID, accounts, buf.Bytes()), nil
}
