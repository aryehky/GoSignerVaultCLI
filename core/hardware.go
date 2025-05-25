package core

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/ledger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// HardwareWallet represents a connected hardware wallet device
type HardwareWallet struct {
	device accounts.Wallet
	path   accounts.DerivationPath
}

// NewHardwareWallet initializes a new hardware wallet connection
func NewHardwareWallet() (*HardwareWallet, error) {
	hub, err := ledger.NewLedgerHub()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ledger hub: %v", err)
	}

	wallets := hub.Wallets()
	if len(wallets) == 0 {
		return nil, errors.New("no hardware wallet found")
	}

	// Use the first available wallet
	wallet := wallets[0]

	// Default to first account
	path := accounts.DefaultBaseDerivationPath

	return &HardwareWallet{
		device: wallet,
		path:   path,
	}, nil
}

// GetAddress returns the Ethereum address for the current derivation path
func (hw *HardwareWallet) GetAddress() (common.Address, error) {
	account, err := hw.device.Derive(hw.path, true)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to derive account: %v", err)
	}
	return account.Address, nil
}

// SignTransaction signs a transaction using the hardware wallet
func (hw *HardwareWallet) SignTransaction(tx *Transaction) ([]byte, error) {
	account, err := hw.device.Derive(hw.path, true)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %v", err)
	}

	// Convert transaction to RLP format
	rlpTx, err := tx.ToRLP()
	if err != nil {
		return nil, fmt.Errorf("failed to encode transaction: %v", err)
	}

	// Sign the transaction
	signature, err := hw.device.SignTx(account, tx.ToEthereumTx(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	return signature, nil
}

// SignMessage signs an arbitrary message using the hardware wallet
func (hw *HardwareWallet) SignMessage(message []byte) ([]byte, error) {
	account, err := hw.device.Derive(hw.path, true)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %v", err)
	}

	// Hash the message according to EIP-191
	hash := crypto.Keccak256Hash(message)

	// Sign the hash
	signature, err := hw.device.SignText(account, hash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	return signature, nil
}
