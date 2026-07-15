// Package onchain provides data structures, instruction builders, and quote simulators
// for interacting directly with the Meteora DLMM smart contracts on the Solana blockchain.
//
// Features include:
//   - Account deserialization matching the Anchor IDL (LBPair, BinArray, etc.).
//   - High-precision math and fee calculations mimicking the exact TypeScript SDK behavior.
//   - Instruction builders for Swap, ClaimFee, and ClaimReward.
//   - PDA derivations for all DLMM accounts (BinArray, Oracle, Reserve, etc.).
package onchain
