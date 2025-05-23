package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultKeystoreDir is the default directory for storing keystore files
	DefaultKeystoreDir = ".keystore"
)

// Manager handles keystore operations
type Manager struct {
	keystoreDir string
}

// NewManager creates a new keystore manager
func NewManager(keystoreDir string) (*Manager, error) {
	if keystoreDir == "" {
		keystoreDir = DefaultKeystoreDir
	}

	// Create keystore directory if it doesn't exist
	if err := os.MkdirAll(keystoreDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create keystore directory: %v", err)
	}

	return &Manager{
		keystoreDir: keystoreDir,
	}, nil
}

// SaveKey saves an encrypted key to the keystore
func (m *Manager) SaveKey(key *EncryptedKey, name string) error {
	// Create the keystore file path
	filePath := filepath.Join(m.keystoreDir, fmt.Sprintf("%s.json", name))

	// Marshal the key to JSON
	data, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal key: %v", err)
	}

	// Write the file with restricted permissions
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write keystore file: %v", err)
	}

	return nil
}

// LoadKey loads an encrypted key from the keystore
func (m *Manager) LoadKey(name string) (*EncryptedKey, error) {
	// Create the keystore file path
	filePath := filepath.Join(m.keystoreDir, fmt.Sprintf("%s.json", name))

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %v", err)
	}

	// Unmarshal the key
	var key EncryptedKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key: %v", err)
	}

	return &key, nil
}

// ListKeys returns a list of all keys in the keystore
func (m *Manager) ListKeys() ([]string, error) {
	files, err := os.ReadDir(m.keystoreDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore directory: %v", err)
	}

	var keys []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			keys = append(keys, file.Name()[:len(file.Name())-5])
		}
	}

	return keys, nil
}

// DeleteKey removes a key from the keystore
func (m *Manager) DeleteKey(name string) error {
	filePath := filepath.Join(m.keystoreDir, fmt.Sprintf("%s.json", name))
	return os.Remove(filePath)
} 