package core

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Wallet represents an Ethereum wallet with its private and public keys
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
}

// NewWallet generates a new Ethereum wallet
func NewWallet() (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKeyECDSA,
		Address:    address,
	}, nil
}

// GetPrivateKeyHex returns the private key as a hex string
func (w *Wallet) GetPrivateKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSA(w.PrivateKey))
}

// GetPublicKeyHex returns the public key as a hex string
func (w *Wallet) GetPublicKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSAPub(w.PublicKey))
}

// GetAddress returns the Ethereum address as a hex string
func (w *Wallet) GetAddress() string {
	return w.Address.Hex()
} 