package tx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TransactionRecord represents a historical transaction record
type TransactionRecord struct {
	Hash        common.Hash `json:"hash"`
	From        string      `json:"from"`
	To          string      `json:"to"`
	Value       string      `json:"value"`
	GasUsed     uint64      `json:"gasUsed"`
	GasPrice    string      `json:"gasPrice"`
	BlockNumber uint64      `json:"blockNumber"`
	Status      string      `json:"status"`
	Timestamp   time.Time   `json:"timestamp"`
	Data        string      `json:"data,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// History manages transaction history
type History struct {
	client   *ethclient.Client
	records  map[common.Hash]*TransactionRecord
	mu       sync.RWMutex
	filePath string
}

// NewHistory creates a new transaction history manager
func NewHistory(rpcURL, filePath string) (*History, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %v", err)
	}

	history := &History{
		client:   client,
		records:  make(map[common.Hash]*TransactionRecord),
		filePath: filePath,
	}

	// Load existing history
	if err := history.load(); err != nil {
		return nil, fmt.Errorf("failed to load history: %v", err)
	}

	return history, nil
}

// AddTransaction adds a transaction to the history
func (h *History) AddTransaction(ctx context.Context, hash common.Hash) error {
	// Get transaction details
	tx, isPending, err := h.client.TransactionByHash(ctx, hash)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	// Get receipt if transaction is not pending
	var receipt *types.Receipt
	if !isPending {
		receipt, err = h.client.TransactionReceipt(ctx, hash)
		if err != nil {
			return fmt.Errorf("failed to get receipt: %v", err)
		}
	}

	// Create record
	record := &TransactionRecord{
		Hash:      hash,
		From:      tx.From().String(),
		To:        tx.To().String(),
		Value:     tx.Value().String(),
		GasPrice:  tx.GasPrice().String(),
		Timestamp: time.Now(),
		Data:      fmt.Sprintf("0x%x", tx.Data()),
	}

	if receipt != nil {
		record.GasUsed = receipt.GasUsed
		record.BlockNumber = receipt.BlockNumber.Uint64()
		if receipt.Status == types.ReceiptStatusFailed {
			record.Status = "failed"
		} else {
			record.Status = "success"
		}
	} else {
		record.Status = "pending"
	}

	// Save record
	h.mu.Lock()
	h.records[hash] = record
	h.mu.Unlock()

	// Save to file
	return h.save()
}

// GetTransaction returns a transaction record
func (h *History) GetTransaction(hash common.Hash) (*TransactionRecord, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if record, exists := h.records[hash]; exists {
		return record, nil
	}
	return nil, fmt.Errorf("transaction not found in history")
}

// GetTransactionsByAddress returns all transactions for an address
func (h *History) GetTransactionsByAddress(address string) []*TransactionRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var records []*TransactionRecord
	for _, record := range h.records {
		if record.From == address || record.To == address {
			records = append(records, record)
		}
	}
	return records
}

// GetRecentTransactions returns the most recent transactions
func (h *History) GetRecentTransactions(limit int) []*TransactionRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var records []*TransactionRecord
	for _, record := range h.records {
		records = append(records, record)
	}

	// Sort by timestamp
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.After(records[j].Timestamp)
	})

	if limit > 0 && limit < len(records) {
		records = records[:limit]
	}

	return records
}

// load loads the transaction history from file
func (h *History) load() error {
	if _, err := os.Stat(h.filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return fmt.Errorf("failed to read history file: %v", err)
	}

	var records map[common.Hash]*TransactionRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse history: %v", err)
	}

	h.mu.Lock()
	h.records = records
	h.mu.Unlock()

	return nil
}

// save saves the transaction history to file
func (h *History) save() error {
	h.mu.RLock()
	data, err := json.MarshalIndent(h.records, "", "  ")
	h.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(h.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(h.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write history file: %v", err)
	}

	return nil
}

// Close closes the history manager
func (h *History) Close() {
	if h.client != nil {
		h.client.Close()
	}
}
