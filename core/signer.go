package core

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// Transaction represents an Ethereum transaction
type Transaction struct {
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	To       *common.Address
	Value    *big.Int
	Data     []byte
	ChainID  *big.Int
}

// SignTransaction signs a transaction with the given private key
func SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) (string, error) {
	// Create the transaction
	ethereumTx := types.NewTransaction(
		tx.Nonce,
		*tx.To,
		tx.Value,
		tx.GasLimit,
		tx.GasPrice,
		tx.Data,
	)

	// Sign the transaction
	signedTx, err := types.SignTx(ethereumTx, types.NewEIP155Signer(tx.ChainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Encode the transaction
	rawTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to encode transaction: %v", err)
	}

	return fmt.Sprintf("0x%x", rawTx), nil
}

// SignMessage signs a message using EIP-191
func SignMessage(message []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	// Create the message hash
	hash := crypto.Keccak256Hash(message)

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	return fmt.Sprintf("0x%x", signature), nil
}

// VerifyMessage verifies a signed message
func VerifyMessage(message []byte, signature string, address common.Address) (bool, error) {
	// Decode the signature
	sig, err := hex.DecodeString(signature[2:]) // Remove "0x" prefix
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Create the message hash
	hash := crypto.Keccak256Hash(message)

	// Recover the public key
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %v", err)
	}

	// Get the address from the public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Compare addresses
	return recoveredAddr == address, nil
}
