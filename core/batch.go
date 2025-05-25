package core

import (
	"encoding/json"
	"fmt"
	"sync"
)

// BatchSigner handles signing multiple transactions in parallel
type BatchSigner struct {
	wallet *Wallet
}

// NewBatchSigner creates a new batch signer
func NewBatchSigner(wallet *Wallet) *BatchSigner {
	return &BatchSigner{
		wallet: wallet,
	}
}

// BatchSignResult represents the result of a batch signing operation
type BatchSignResult struct {
	TransactionID string `json:"transactionId"`
	Signature     []byte `json:"signature"`
	Error         string `json:"error,omitempty"`
}

// SignBatch signs multiple transactions in parallel
func (bs *BatchSigner) SignBatch(transactions []*Transaction) []BatchSignResult {
	var wg sync.WaitGroup
	results := make([]BatchSignResult, len(transactions))

	// Create a channel to collect results
	resultChan := make(chan struct {
		index  int
		result BatchSignResult
	}, len(transactions))

	// Sign each transaction in a goroutine
	for i, tx := range transactions {
		wg.Add(1)
		go func(index int, transaction *Transaction) {
			defer wg.Done()

			result := BatchSignResult{
				TransactionID: fmt.Sprintf("tx_%d", index),
			}

			// Sign the transaction
			signature, err := bs.wallet.SignTransaction(transaction)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Signature = signature
			}

			resultChan <- struct {
				index  int
				result BatchSignResult
			}{index, result}
		}(i, tx)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		results[result.index] = result.result
	}

	return results
}

// BatchSignResultToJSON converts a batch sign result to JSON
func BatchSignResultToJSON(results []BatchSignResult) (string, error) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %v", err)
	}
	return string(data), nil
}

// BatchSignResultFromJSON parses a JSON string into batch sign results
func BatchSignResultFromJSON(jsonData string) ([]BatchSignResult, error) {
	var results []BatchSignResult
	if err := json.Unmarshal([]byte(jsonData), &results); err != nil {
		return nil, fmt.Errorf("failed to parse results: %v", err)
	}
	return results, nil
}
