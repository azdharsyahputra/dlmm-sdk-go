package onchain

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Client wraps a Solana RPC client to interact with the DLMM program.
type Client struct {
	rpcClient *rpc.Client
}

// NewClient returns a new DLMM On-Chain Client.
func NewClient(endpoint string) *Client {
	if endpoint == "" {
		endpoint = rpc.MainNetBeta_RPC
	}
	return &Client{
		rpcClient: rpc.New(endpoint),
	}
}

// GetRPCClient returns the underlying Solana RPC client.
func (c *Client) GetRPCClient() *rpc.Client {
	return c.rpcClient
}

// GetLBPair fetches and deserializes the state of a DLMM pool (LBPair).
func (c *Client) GetLBPair(ctx context.Context, address solana.PublicKey) (*LBPair, error) {
	resp, err := c.rpcClient.GetAccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get LBPair account info: %w", err)
	}
	if resp == nil || resp.Value == nil {
		return nil, fmt.Errorf("lbPair account not found")
	}

	var pair LBPair
	if err := DecodeLBPair(resp.Value.Data.GetBinary(), &pair); err != nil {
		return nil, err
	}
	return &pair, nil
}

// GetBinArray fetches and deserializes a BinArray account.
func (c *Client) GetBinArray(ctx context.Context, address solana.PublicKey) (*BinArray, error) {
	resp, err := c.rpcClient.GetAccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get BinArray account info: %w", err)
	}
	if resp == nil || resp.Value == nil {
		return nil, fmt.Errorf("BinArray account not found")
	}

	var binArray BinArray
	if err := DecodeLBPair(resp.Value.Data.GetBinary(), &binArray); err != nil {
		return nil, err
	}
	return &binArray, nil
}

// GetBinArrayByIndex derives, fetches, and deserializes a BinArray by its index.
func (c *Client) GetBinArrayByIndex(ctx context.Context, lbPair solana.PublicKey, index int64) (*BinArray, error) {
	pubkey, _, err := DeriveBinArray(lbPair, index)
	if err != nil {
		return nil, fmt.Errorf("failed to derive BinArray address: %w", err)
	}
	return c.GetBinArray(ctx, pubkey)
}

// GetBin extracts a specific bin by bin ID from a deserialized BinArray.
func (binArray *BinArray) GetBin(binID int32) (*Bin, error) {
	lowerBinID, upperBinID := GetBinArrayLowerUpperBinID(binArray.Index)
	if binID < lowerBinID || binID > upperBinID {
		return nil, fmt.Errorf("binID %d is out of range for BinArray index %d [%d, %d]", binID, binArray.Index, lowerBinID, upperBinID)
	}
	idx := binID - lowerBinID
	return &binArray.Bins[idx], nil
}

// GetPositionV2 fetches and deserializes a user position.
func (c *Client) GetPositionV2(ctx context.Context, address solana.PublicKey) (*PositionV2, error) {
	resp, err := c.rpcClient.GetAccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get PositionV2 account info: %w", err)
	}
	if resp == nil || resp.Value == nil {
		return nil, fmt.Errorf("PositionV2 account not found")
	}

	var pos PositionV2
	if err := DecodeLBPair(resp.Value.Data.GetBinary(), &pos); err != nil {
		return nil, err
	}
	return &pos, nil
}

// GetBinArrayBitmapExtension fetches and deserializes the bitmap extension account.
func (c *Client) GetBinArrayBitmapExtension(ctx context.Context, address solana.PublicKey) (*BinArrayBitmapExtension, error) {
	resp, err := c.rpcClient.GetAccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get BinArrayBitmapExtension account info: %w", err)
	}
	if resp == nil || resp.Value == nil {
		return nil, fmt.Errorf("BinArrayBitmapExtension account not found")
	}

	var ext BinArrayBitmapExtension
	if err := DecodeLBPair(resp.Value.Data.GetBinary(), &ext); err != nil {
		return nil, err
	}
	return &ext, nil
}
