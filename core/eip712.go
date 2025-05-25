package core

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// TypedData represents EIP-712 typed data structure
type TypedData struct {
	Types       apitypes.Types           `json:"types"`
	PrimaryType string                   `json:"primaryType"`
	Domain      apitypes.TypedDataDomain `json:"domain"`
	Message     map[string]interface{}   `json:"message"`
}

// SignTypedData signs an EIP-712 typed data message
func (w *Wallet) SignTypedData(data *TypedData) ([]byte, error) {
	// Convert to Ethereum's internal format
	typedData := apitypes.TypedData{
		Types:       data.Types,
		PrimaryType: data.PrimaryType,
		Domain:      data.Domain,
		Message:     data.Message,
	}

	// Get the domain separator
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to hash domain separator: %v", err)
	}

	// Get the message hash
	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to hash message: %v", err)
	}

	// Create the final hash
	hash := crypto.Keccak256Hash(
		[]byte("\x19\x01"),
		domainSeparator,
		messageHash,
	)

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), w.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign typed data: %v", err)
	}

	return signature, nil
}

// ParseTypedData parses a JSON string into a TypedData structure
func ParseTypedData(jsonData string) (*TypedData, error) {
	var data TypedData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, fmt.Errorf("failed to parse typed data: %v", err)
	}
	return &data, nil
}

// VerifyTypedDataSignature verifies an EIP-712 signature
func VerifyTypedDataSignature(data *TypedData, signature []byte) (common.Address, error) {
	// Convert to Ethereum's internal format
	typedData := apitypes.TypedData{
		Types:       data.Types,
		PrimaryType: data.PrimaryType,
		Domain:      data.Domain,
		Message:     data.Message,
	}

	// Get the domain separator
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to hash domain separator: %v", err)
	}

	// Get the message hash
	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to hash message: %v", err)
	}

	// Create the final hash
	hash := crypto.Keccak256Hash(
		[]byte("\x19\x01"),
		domainSeparator,
		messageHash,
	)

	// Recover the public key
	pubKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover public key: %v", err)
	}

	// Get the address
	address := crypto.PubkeyToAddress(*pubKey)
	return address, nil
}
