package tx

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SimulationResult represents the result of a transaction simulation
type SimulationResult struct {
	Success      bool              `json:"success"`
	GasUsed      uint64            `json:"gasUsed"`
	GasPrice     *big.Int          `json:"gasPrice"`
	TotalCost    *big.Int          `json:"totalCost"`
	Error        string            `json:"error,omitempty"`
	Trace        []string          `json:"trace,omitempty"`
	StateChanges map[string]string `json:"stateChanges,omitempty"`
}

// Simulator handles transaction simulation and gas estimation
type Simulator struct {
	client *ethclient.Client
}

// NewSimulator creates a new transaction simulator
func NewSimulator(rpcURL string) (*Simulator, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %v", err)
	}

	return &Simulator{
		client: client,
	}, nil
}

// EstimateGas estimates the gas required for a transaction
func (s *Simulator) EstimateGas(ctx context.Context, tx *Transaction) (uint64, error) {
	// Convert to Ethereum transaction
	ethTx := tx.ToEthereumTx()

	// Create call message
	msg := ethereum.CallMsg{
		From:     ethTx.From(),
		To:       ethTx.To(),
		Gas:      ethTx.Gas(),
		GasPrice: ethTx.GasPrice(),
		Value:    ethTx.Value(),
		Data:     ethTx.Data(),
	}

	// Estimate gas
	gasLimit, err := s.client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %v", err)
	}

	return gasLimit, nil
}

// SimulateTransaction simulates a transaction and returns detailed results
func (s *Simulator) SimulateTransaction(ctx context.Context, tx *Transaction) (*SimulationResult, error) {
	// Convert to Ethereum transaction
	ethTx := tx.ToEthereumTx()

	// Create call message
	msg := ethereum.CallMsg{
		From:     ethTx.From(),
		To:       ethTx.To(),
		Gas:      ethTx.Gas(),
		GasPrice: ethTx.GasPrice(),
		Value:    ethTx.Value(),
		Data:     ethTx.Data(),
	}

	// Get current block number
	blockNumber, err := s.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get block number: %v", err)
	}

	// Simulate transaction
	result := &SimulationResult{
		StateChanges: make(map[string]string),
	}

	// Call the transaction
	_, err = s.client.CallContract(ctx, msg, big.NewInt(int64(blockNumber)))
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	// Get gas price
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	// Estimate gas
	gasLimit, err := s.client.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %v", err)
	}

	// Calculate total cost
	totalCost := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasLimit))

	result.Success = true
	result.GasUsed = gasLimit
	result.GasPrice = gasPrice
	result.TotalCost = totalCost

	return result, nil
}

// GetGasPrice returns the current gas price
func (s *Simulator) GetGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}
	return gasPrice, nil
}

// GetGasPriceHistory returns historical gas prices
func (s *Simulator) GetGasPriceHistory(ctx context.Context, blocks int) ([]*big.Int, error) {
	// Get current block number
	currentBlock, err := s.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get block number: %v", err)
	}

	var prices []*big.Int
	for i := 0; i < blocks; i++ {
		blockNumber := currentBlock - uint64(i)
		block, err := s.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			return nil, fmt.Errorf("failed to get block %d: %v", blockNumber, err)
		}

		// Get base fee if available (EIP-1559)
		if block.BaseFee() != nil {
			prices = append(prices, block.BaseFee())
		} else {
			// Fallback to gas price from the first transaction
			if len(block.Transactions()) > 0 {
				prices = append(prices, block.Transactions()[0].GasPrice())
			}
		}
	}

	return prices, nil
}

// Close closes the RPC connection
func (s *Simulator) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
