package unit_test

import (
	"testing"

	"github.com/azdharsyahputra/dlmm-sdk-go/onchain"
	"github.com/shopspring/decimal"
)

func TestDecimalPriceMath(t *testing.T) {
	binID := int32(-6388)
	binStep := uint16(4)
	decimalsX := int32(9)
	decimalsY := int32(6)

	// 1. Convert bin to lamport price
	priceDec := onchain.BinIDToPriceDecimal(binID, binStep)
	priceFloat := onchain.BinIDToPrice(binID, binStep)

	priceDecFloat, _ := priceDec.Float64()
	if priceDecFloat != priceFloat {
		t.Errorf("Mismatch in BinIDToPrice: decimal=%f, float=%f", priceDecFloat, priceFloat)
	}

	// 2. Convert bin to token price
	tokenPriceDec := onchain.BinIDToTokenPriceDecimal(binID, binStep, decimalsX, decimalsY)
	tokenPriceFloat := onchain.BinIDToTokenPrice(binID, binStep, decimalsX, decimalsY)

	tokenPriceDecFloat, _ := tokenPriceDec.Float64()
	if tokenPriceDecFloat != tokenPriceFloat {
		t.Errorf("Mismatch in BinIDToTokenPrice: decimal=%f, float=%f", tokenPriceDecFloat, tokenPriceFloat)
	}

	// 3. Convert token price to bin ID
	targetPrice := decimal.NewFromFloat(77.716391)
	calculatedBinID := onchain.TokenPriceToBinIDDecimal(targetPrice, binStep, decimalsX, decimalsY)

	if calculatedBinID != binID {
		t.Errorf("Expected bin ID %d, got %d for token price %s", binID, calculatedBinID, targetPrice.String())
	}

	t.Logf("Decimal Price Math verified successfully!")
	t.Logf("Token Price Decimal: %s", tokenPriceDec.String())
}
