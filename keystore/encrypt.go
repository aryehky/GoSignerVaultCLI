package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
)

// EncryptedKey represents an encrypted private key
type EncryptedKey struct {
	Address    string          `json:"address"`
	Crypto     CryptoJSON     `json:"crypto"`
	Version    int            `json:"version"`
	ID         string         `json:"id"`
}

// CryptoJSON represents the encrypted data structure
type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams CipherParamsJSON      `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

// CipherParamsJSON represents the cipher parameters
type CipherParamsJSON struct {
	IV string `json:"iv"`
}

// EncryptKey encrypts a private key using AES-256-GCM
func EncryptKey(privateKey []byte, password string) (*EncryptedKey, error) {
	// Generate a random salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	// Derive key from password
	derivedKey := deriveKey(password, salt)

	// Generate random IV
	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Encrypt the private key
	ciphertext := aesGCM.Seal(nil, iv, privateKey, nil)

	// Create MAC
	mac := crypto.Keccak256(append(derivedKey[16:32], ciphertext...))

	// Create the encrypted key structure
	encryptedKey := &EncryptedKey{
		Address: crypto.PubkeyToAddress(crypto.ToECDSA(privateKey).PublicKey).Hex(),
		Crypto: CryptoJSON{
			Cipher:     "aes-256-gcm",
			CipherText: fmt.Sprintf("0x%x", ciphertext),
			CipherParams: CipherParamsJSON{
				IV: fmt.Sprintf("0x%x", iv),
			},
			KDF: "pbkdf2",
			KDFParams: map[string]interface{}{
				"c":     262144,
				"dklen": 32,
				"prf":   "hmac-sha256",
				"salt":  fmt.Sprintf("0x%x", salt),
			},
			MAC: fmt.Sprintf("0x%x", mac),
		},
		Version: 3,
		ID:      fmt.Sprintf("%x", crypto.Keccak256([]byte("GoSignerVaultCLI"))),
	}

	return encryptedKey, nil
}

// deriveKey derives an encryption key from a password and salt
func deriveKey(password string, salt []byte) []byte {
	// Simple key derivation using SHA256
	// In production, use a proper KDF like PBKDF2
	key := sha256.Sum256(append([]byte(password), salt...))
	return key[:]
} 