import os
import re

# 1. client package
client_dir = "/Users/csadeveloper/go-project/sdk-go/client"
for fname in os.listdir(client_dir):
    if not fname.endswith(".go"): continue
    path = os.path.join(client_dir, fname)
    with open(path, "r") as f: content = f.read()
    
    # Static fmt.Errorf to errors.New
    content = re.sub(r'fmt\.Errorf\("([^%]+?)"\)', r'errors.New("\1")', content)
    
    # Ensure errors is imported if errors.New is used
    if 'errors.New(' in content and '"errors"' not in content:
        content = re.sub(r'import \(', 'import (\n\t"errors"', content, count=1)
        if 'import (' not in content: # single import
            content = re.sub(r'import "fmt"', 'import (\n\t"errors"\n\t"fmt"\n)', content)
            
    # Rename Option to ClientOption
    content = content.replace("type Option func(*Client)", "type ClientOption func(*Client)")
    content = content.replace(") Option {", ") ClientOption {")
    content = content.replace("...Option)", "...ClientOption)")
    
    # Exported Names
    content = content.replace("PaginationResponse_OpenLimitOrderPoolSummary", "OpenLimitOrderPoolSummaryPage")
    content = content.replace("PaginationResponse_ClosedLimitOrderPoolSummary", "ClosedLimitOrderPoolSummaryPage")
    
    # Fix GetPortfolio (Required param)
    if fname == "portfolio.go":
        content = content.replace("func (c *Client) GetPortfolio(ctx context.Context, params *GetPortfolioParams) (*PortfolioResponse, error) {",
                                  "func (c *Client) GetPortfolio(ctx context.Context, user string, params *GetPortfolioParams) (*PortfolioResponse, error) {")
        content = content.replace('url := fmt.Sprintf("%s/portfolio/%s", c.baseURL, params.User)',
                                  'url := fmt.Sprintf("%s/portfolio/%s", c.baseURL, user)')
                                  
    with open(path, "w") as f: f.write(content)

# 2. onchain package
onchain_dir = "/Users/csadeveloper/go-project/sdk-go/onchain"
for fname in os.listdir(onchain_dir):
    if not fname.endswith(".go"): continue
    path = os.path.join(onchain_dir, fname)
    with open(path, "r") as f: content = f.read()
    
    # LbPair -> LBPair
    content = content.replace("LbPair", "LBPair")
    content = content.replace("lbPair", "lbPair") # keep unexported as lbPair
    
    # Id -> ID
    content = content.replace("ActiveId", "ActiveID")
    content = content.replace("BinId", "BinID")
    content = content.replace("binId", "binID")
    
    # ILM_BASE -> ilmBase
    content = content.replace("ILM_BASE", "ilmBase")
    
    # Error capitalization
    content = content.replace('fmt.Errorf("LbPair', 'fmt.Errorf("lbPair')
    content = content.replace('fmt.Errorf("LBPair', 'fmt.Errorf("lbPair')
    
    # BinArray allocations
    content = content.replace("[]BinArray", "[]*BinArray")
    content = content.replace("map[int64]BinArray", "map[int64]*BinArray")
    
    # Helper func removal
    content = content.replace("GetBinFromBinArrayHelper", "GetBin")
    content = content.replace("func GetBin(binID int32, binArray *BinArray)", "func (binArray *BinArray) GetBin(binID int32)")
    
    # Bin pointer
    content = content.replace("func SwapExactInQuoteAtBin(bin Bin", "func swapExactInQuoteAtBin(bin *Bin")
    
    # DeserializeAccount unexported
    content = content.replace("func DeserializeAccount(", "func deserializeAccount(")
    
    with open(path, "w") as f: f.write(content)
