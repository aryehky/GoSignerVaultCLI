package keystore

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupConfig represents the configuration for a keystore backup
type BackupConfig struct {
	Version   string            `json:"version"`
	Timestamp int64             `json:"timestamp"`
	Keystores []string          `json:"keystores"`
	Metadata  map[string]string `json:"metadata"`
}

// CreateBackup creates an encrypted backup of the keystore directory
func CreateBackup(keystoreDir string, backupPath string, password string) error {
	// Create a temporary directory for the backup
	tempDir, err := os.MkdirTemp("", "keystore-backup-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create backup config
	config := BackupConfig{
		Version:   "1.0",
		Timestamp: time.Now().Unix(),
		Metadata:  make(map[string]string),
	}

	// Copy keystore files to temp directory
	err = filepath.Walk(keystoreDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Skip non-keystore files
		if filepath.Ext(path) != ".json" {
			return nil
		}

		// Copy file to temp directory
		relPath, err := filepath.Rel(keystoreDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(tempDir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
			return err
		}

		if err := copyFile(path, destPath); err != nil {
			return err
		}

		config.Keystores = append(config.Keystores, relPath)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy keystore files: %v", err)
	}

	// Create config file
	configPath := filepath.Join(tempDir, "backup.json")
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, configData, 0600); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	// Create encrypted zip archive
	if err := createEncryptedZip(tempDir, backupPath, password); err != nil {
		return fmt.Errorf("failed to create encrypted backup: %v", err)
	}

	return nil
}

// RestoreBackup restores a keystore backup to the specified directory
func RestoreBackup(backupPath string, keystoreDir string, password string) error {
	// Create temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "keystore-restore-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract encrypted zip
	if err := extractEncryptedZip(backupPath, tempDir, password); err != nil {
		return fmt.Errorf("failed to extract backup: %v", err)
	}

	// Read backup config
	configPath := filepath.Join(tempDir, "backup.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	var config BackupConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	// Restore keystore files
	for _, keystorePath := range config.Keystores {
		srcPath := filepath.Join(tempDir, keystorePath)
		destPath := filepath.Join(keystoreDir, keystorePath)

		if err := os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		if err := copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy keystore file: %v", err)
		}
	}

	return nil
}

// Helper function to copy a file
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// Helper function to create an encrypted zip archive
func createEncryptedZip(srcDir, zipPath, password string) error {
	// Create zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through source directory
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Create zip file entry
		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Read and encrypt file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		encryptedData, err := encryptData(data, password)
		if err != nil {
			return err
		}

		_, err = writer.Write(encryptedData)
		return err
	})

	return err
}

// Helper function to extract an encrypted zip archive
func extractEncryptedZip(zipPath, destDir, password string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Create destination file
		destPath := filepath.Join(destDir, file.Name)
		if err := os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
			return err
		}

		// Open source file
		rc, err := file.Open()
		if err != nil {
			return err
		}

		// Read and decrypt data
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}

		decryptedData, err := decryptData(data, password)
		if err != nil {
			return err
		}

		// Write decrypted data
		if err := os.WriteFile(destPath, decryptedData, 0600); err != nil {
			return err
		}
	}

	return nil
}

// Helper function to encrypt data with AES-256-GCM
func encryptData(data []byte, password string) ([]byte, error) {
	// Derive key from password
	key := deriveKey(password)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Helper function to decrypt data with AES-256-GCM
func decryptData(data []byte, password string) ([]byte, error) {
	// Derive key from password
	key := deriveKey(password)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Helper function to derive a key from a password
func deriveKey(password string) []byte {
	// In a real implementation, use a proper key derivation function like PBKDF2
	// This is a simplified version for demonstration
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}
