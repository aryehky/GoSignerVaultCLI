package cmd

import (
	"fmt"

	"github.com/aryehky/gosignervaultcli/core"
	"github.com/aryehky/gosignervaultcli/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var (
	keystoreDir string
	keyName     string
	password    string
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage wallet keys",
	Long:  `Generate, list, and manage wallet keys for signing transactions.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new wallet key",
	Long:  `Generate a new Ethereum wallet key and save it to the keystore.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create keystore manager
		manager, err := keystore.NewManager(keystoreDir)
		if err != nil {
			return fmt.Errorf("failed to create keystore manager: %v", err)
		}

		// Generate new wallet
		wallet, err := core.NewWallet()
		if err != nil {
			return fmt.Errorf("failed to generate wallet: %v", err)
		}

		// Encrypt private key
		encryptedKey, err := keystore.EncryptKey(crypto.FromECDSA(wallet.PrivateKey), password)
		if err != nil {
			return fmt.Errorf("failed to encrypt key: %v", err)
		}

		// Save to keystore
		if err := manager.SaveKey(encryptedKey, keyName); err != nil {
			return fmt.Errorf("failed to save key: %v", err)
		}

		fmt.Printf("Generated new wallet: %s\n", wallet.GetAddress())
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all wallet keys",
	Long:  `List all wallet keys stored in the keystore.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create keystore manager
		manager, err := keystore.NewManager(keystoreDir)
		if err != nil {
			return fmt.Errorf("failed to create keystore manager: %v", err)
		}

		// List keys
		keys, err := manager.ListKeys()
		if err != nil {
			return fmt.Errorf("failed to list keys: %v", err)
		}

		if len(keys) == 0 {
			fmt.Println("No keys found in keystore")
			return nil
		}

		fmt.Println("Available keys:")
		for _, key := range keys {
			fmt.Printf("- %s\n", key)
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a wallet key",
	Long:  `Delete a wallet key from the keystore.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create keystore manager
		manager, err := keystore.NewManager(keystoreDir)
		if err != nil {
			return fmt.Errorf("failed to create keystore manager: %v", err)
		}

		// Delete key
		if err := manager.DeleteKey(keyName); err != nil {
			return fmt.Errorf("failed to delete key: %v", err)
		}

		fmt.Printf("Deleted key: %s\n", keyName)
		return nil
	},
}

func init() {
	// Add flags
	keysCmd.PersistentFlags().StringVar(&keystoreDir, "keystore", ".keystore", "Keystore directory")
	generateCmd.Flags().StringVar(&keyName, "name", "", "Key name")
	generateCmd.Flags().StringVar(&password, "password", "", "Encryption password")
	deleteCmd.Flags().StringVar(&keyName, "name", "", "Key name to delete")

	// Mark required flags
	generateCmd.MarkFlagRequired("name")
	generateCmd.MarkFlagRequired("password")
	deleteCmd.MarkFlagRequired("name")

	// Add commands
	keysCmd.AddCommand(generateCmd)
	keysCmd.AddCommand(listCmd)
	keysCmd.AddCommand(deleteCmd)
}
