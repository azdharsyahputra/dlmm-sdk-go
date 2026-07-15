// Package dlmm provides the official Go SDK for Meteora's Dynamic Liquidity Market Maker (DLMM).
//
// This SDK is divided into two primary sub-packages tailored for different use cases:
//
// 1. client: A REST API client for interacting with Meteora's off-chain analytics,
// fetching pools, historical data (OHLCV), volume, and user portfolios.
//
// 2. onchain: Core blockchain utilities for parsing raw Solana accounts (LBPairs, BinArrays),
// calculating exact swap quotes locally, and constructing instructions for Swaps, 
// Limit Orders, and Fee Claiming.
//
// # Examples
//
// Full executable examples can be found in the "examples/" directory of this repository,
// including a complete Devnet E2E swap simulation.
//
// # Sub-packages
//
// Please navigate to the "client" or "onchain" packages below to see their
// respective documentation, types, and runnable examples.
package dlmm
