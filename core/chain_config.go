package core

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
)

// ChainConfig represents the configuration for an EVM-compatible chain
type ChainConfig struct {
	Name      string   `json:"name"`
	ChainID   *big.Int `json:"chainId"`
	RPCURL    string   `json:"rpcUrl"`
	Symbol    string   `json:"symbol"`
	Explorer  string   `json:"explorer"`
	IsTestnet bool     `json:"isTestnet"`
}

// DefaultChains contains predefined chain configurations
var DefaultChains = map[string]*ChainConfig{
	"ethereum": {
		Name:      "Ethereum Mainnet",
		ChainID:   big.NewInt(1),
		RPCURL:    "https://mainnet.infura.io/v3/YOUR-PROJECT-ID",
		Symbol:    "ETH",
		Explorer:  "https://etherscan.io",
		IsTestnet: false,
	},
	"polygon": {
		Name:      "Polygon Mainnet",
		ChainID:   big.NewInt(137),
		RPCURL:    "https://polygon-rpc.com",
		Symbol:    "MATIC",
		Explorer:  "https://polygonscan.com",
		IsTestnet: false,
	},
	"bsc": {
		Name:      "BNB Smart Chain",
		ChainID:   big.NewInt(56),
		RPCURL:    "https://bsc-dataseed.binance.org",
		Symbol:    "BNB",
		Explorer:  "https://bscscan.com",
		IsTestnet: false,
	},
	"avalanche": {
		Name:      "Avalanche C-Chain",
		ChainID:   big.NewInt(43114),
		RPCURL:    "https://api.avax.network/ext/bc/C/rpc",
		Symbol:    "AVAX",
		Explorer:  "https://snowtrace.io",
		IsTestnet: false,
	},
}

// LoadChainConfig loads chain configurations from a JSON file
func LoadChainConfig(path string) (map[string]*ChainConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read chain config file: %v", err)
	}

	var configs map[string]*ChainConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse chain config file: %v", err)
	}

	return configs, nil
}

// SaveChainConfig saves chain configurations to a JSON file
func SaveChainConfig(path string, configs map[string]*ChainConfig) error {
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chain configs: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write chain config file: %v", err)
	}

	return nil
}

// GetChainConfig returns a chain configuration by name
func GetChainConfig(name string) (*ChainConfig, error) {
	config, ok := DefaultChains[name]
	if !ok {
		return nil, fmt.Errorf("chain %s not found", name)
	}
	return config, nil
}
