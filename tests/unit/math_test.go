package unit_test

import (
	"testing"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/shopspring/decimal"
)

func TestDecimalPriceMath(t *testing.T) {
	binId := int32(-6388)
	binStep := uint16(4)
	decimalsX := int32(9)
	decimalsY := int32(6)

	// 1. Convert bin to lamport price
	priceDec := onchain.BinIdToPriceDecimal(binId, binStep)
	priceFloat := onchain.BinIdToPrice(binId, binStep)

	priceDecFloat, _ := priceDec.Float64()
	if priceDecFloat != priceFloat {
		t.Errorf("Mismatch in BinIdToPrice: decimal=%f, float=%f", priceDecFloat, priceFloat)
	}

	// 2. Convert bin to token price
	tokenPriceDec := onchain.BinIdToTokenPriceDecimal(binId, binStep, decimalsX, decimalsY)
	tokenPriceFloat := onchain.BinIdToTokenPrice(binId, binStep, decimalsX, decimalsY)

	tokenPriceDecFloat, _ := tokenPriceDec.Float64()
	if tokenPriceDecFloat != tokenPriceFloat {
		t.Errorf("Mismatch in BinIdToTokenPrice: decimal=%f, float=%f", tokenPriceDecFloat, tokenPriceFloat)
	}

	// 3. Convert token price to bin ID
	targetPrice := decimal.NewFromFloat(77.716391)
	calculatedBinId := onchain.TokenPriceToBinIdDecimal(targetPrice, binStep, decimalsX, decimalsY)

	if calculatedBinId != binId {
		t.Errorf("Expected bin ID %d, got %d for token price %s", binId, calculatedBinId, targetPrice.String())
	}

	t.Logf("Decimal Price Math verified successfully!")
	t.Logf("Token Price Decimal: %s", tokenPriceDec.String())
}
