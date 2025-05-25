package tx

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TransactionStatus represents the status of a monitored transaction
type TransactionStatus struct {
	Hash      common.Hash `json:"hash"`
	Status    string      `json:"status"`
	BlockNum  uint64      `json:"blockNum,omitempty"`
	GasUsed   uint64      `json:"gasUsed,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Monitor handles transaction monitoring
type Monitor struct {
	client    *ethclient.Client
	statuses  map[common.Hash]*TransactionStatus
	mu        sync.RWMutex
	callbacks map[common.Hash][]func(*TransactionStatus)
}

// NewMonitor creates a new transaction monitor
func NewMonitor(rpcURL string) (*Monitor, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %v", err)
	}

	return &Monitor{
		client:    client,
		statuses:  make(map[common.Hash]*TransactionStatus),
		callbacks: make(map[common.Hash][]func(*TransactionStatus)),
	}, nil
}

// MonitorTransaction starts monitoring a transaction
func (m *Monitor) MonitorTransaction(ctx context.Context, hash common.Hash) error {
	m.mu.Lock()
	if _, exists := m.statuses[hash]; exists {
		m.mu.Unlock()
		return fmt.Errorf("transaction already being monitored")
	}

	status := &TransactionStatus{
		Hash:      hash,
		Status:    "pending",
		Timestamp: time.Now(),
	}
	m.statuses[hash] = status
	m.mu.Unlock()

	// Start monitoring in a goroutine
	go m.monitorTransaction(ctx, hash)

	return nil
}

// monitorTransaction continuously monitors a transaction
func (m *Monitor) monitorTransaction(ctx context.Context, hash common.Hash) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			receipt, err := m.client.TransactionReceipt(ctx, hash)
			if err != nil {
				if err.Error() == "not found" {
					continue
				}
				m.updateStatus(hash, "error", 0, 0, err.Error())
				return
			}

			status := "success"
			if receipt.Status == types.ReceiptStatusFailed {
				status = "failed"
			}

			m.updateStatus(hash, status, receipt.BlockNumber.Uint64(), receipt.GasUsed, "")
			return
		}
	}
}

// updateStatus updates the status of a transaction
func (m *Monitor) updateStatus(hash common.Hash, status string, blockNum, gasUsed uint64, errMsg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if txStatus, exists := m.statuses[hash]; exists {
		txStatus.Status = status
		txStatus.BlockNum = blockNum
		txStatus.GasUsed = gasUsed
		txStatus.Error = errMsg
		txStatus.Timestamp = time.Now()

		// Call callbacks
		if callbacks, exists := m.callbacks[hash]; exists {
			for _, callback := range callbacks {
				callback(txStatus)
			}
		}
	}
}

// GetStatus returns the current status of a transaction
func (m *Monitor) GetStatus(hash common.Hash) (*TransactionStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if status, exists := m.statuses[hash]; exists {
		return status, nil
	}
	return nil, fmt.Errorf("transaction not being monitored")
}

// AddCallback adds a callback function for status updates
func (m *Monitor) AddCallback(hash common.Hash, callback func(*TransactionStatus)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callbacks[hash] = append(m.callbacks[hash], callback)
}

// RemoveCallback removes a callback function
func (m *Monitor) RemoveCallback(hash common.Hash, callback func(*TransactionStatus)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if callbacks, exists := m.callbacks[hash]; exists {
		for i, cb := range callbacks {
			if fmt.Sprintf("%v", cb) == fmt.Sprintf("%v", callback) {
				m.callbacks[hash] = append(callbacks[:i], callbacks[i+1:]...)
				break
			}
		}
	}
}

// Close closes the monitor
func (m *Monitor) Close() {
	if m.client != nil {
		m.client.Close()
	}
}
