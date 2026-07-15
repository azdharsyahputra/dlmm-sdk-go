// Package onchain provides tools for interacting with Meteora's DLMM Solana smart contracts.
//
// It provides functionality to parse on-chain account data (LBPair, BinArray), 
// compute swap quotes locally using exact bin math, and build Solana instructions 
// for interactions like swaps and claiming fees.
//
// # Safety and Math
//
// Math operations inside this package mirror the Rust smart contract logic exactly.
// It utilizes math/big extensively for 128-bit operations. To optimize memory allocations,
// all math operations modify variables in-place, and users should be mindful of 
// pointer copies when interacting with large structures like BinArrays.
//
// # Quick Start
//
//	// Fetch LBPair state
//	var pair onchain.LBPair
//	onchain.DecodeLBPair(rawBytes, &pair)
//
//	// Compute a local swap quote to get exact Out amount
//	quote, err := onchain.ComputeSwapQuote(
//		&pair, binArrays, inAmount, swapForY, false, 50, 0, nil, nil,
//	)
package onchain
