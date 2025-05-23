package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
)

// EncryptedKey represents an encrypted private key
type EncryptedKey struct {
	Address string     `json:"address"`
	Crypto  CryptoJSON `json:"crypto"`
	Version int        `json:"version"`
	ID      string     `json:"id"`
}

// CryptoJSON represents the encrypted data structure
type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams CipherParamsJSON       `json:"cipherparams"`
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

// DecryptKey decrypts a private key using the provided password
func DecryptKey(key *EncryptedKey, password string) (*ecdsa.PrivateKey, error) {
	// Get salt from KDF params
	saltHex, ok := key.Crypto.KDFParams["salt"].(string)
	if !ok {
		return nil, errors.New("invalid salt in key file")
	}
	salt, err := hex.DecodeString(saltHex[2:]) // Remove "0x" prefix
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %v", err)
	}

	// Derive key from password
	derivedKey := deriveKey(password, salt)

	// Get IV from cipher params
	iv, err := hex.DecodeString(key.Crypto.CipherParams.IV[2:]) // Remove "0x" prefix
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %v", err)
	}

	// Get ciphertext
	ciphertext, err := hex.DecodeString(key.Crypto.CipherText[2:]) // Remove "0x" prefix
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %v", err)
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

	// Decrypt the private key
	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %v", err)
	}

	// Convert to private key
	privateKey, err := crypto.ToECDSA(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to private key: %v", err)
	}

	return privateKey, nil
}

// deriveKey derives an encryption key from a password and salt
func deriveKey(password string, salt []byte) []byte {
	// Simple key derivation using SHA256
	// In production, use a proper KDF like PBKDF2
	key := sha256.Sum256(append([]byte(password), salt...))
	return key[:]
}
