package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aryehky/gosignervaultcli/core"
	"github.com/aryehky/gosignervaultcli/keystore"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
	chainName  string
	message    string
)

// SignCmd is the root command for signing operations
var SignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign transactions and messages",
	Long:  `Sign Ethereum transactions and messages using stored wallet keys.`,
}

var signTxCmd = &cobra.Command{
	Use:   "tx",
	Short: "Sign a transaction",
	Long:  `Sign an Ethereum transaction using a stored wallet key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load chain config
		chain, err := core.GetChainConfig(chainName)
		if err != nil {
			return fmt.Errorf("failed to get chain config: %v", err)
		}

		// Read input file
		data, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %v", err)
		}

		// Parse transaction
		var tx core.Transaction
		if err := json.Unmarshal(data, &tx); err != nil {
			return fmt.Errorf("failed to parse transaction: %v", err)
		}

		// Set chain ID
		tx.ChainID = chain.ChainID

		// Load key
		manager, err := keystore.NewManager(keystoreDir)
		if err != nil {
			return fmt.Errorf("failed to create keystore manager: %v", err)
		}

		encryptedKey, err := manager.LoadKey(keyName)
		if err != nil {
			return fmt.Errorf("failed to load key: %v", err)
		}

		// Decrypt key
		privateKey, err := keystore.DecryptKey(encryptedKey, password)
		if err != nil {
			return fmt.Errorf("failed to decrypt key: %v", err)
		}

		// Sign transaction
		signedTx, err := core.SignTransaction(&tx, privateKey)
		if err != nil {
			return fmt.Errorf("failed to sign transaction: %v", err)
		}

		// Write output
		if err := ioutil.WriteFile(outputFile, []byte(signedTx), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %v", err)
		}

		fmt.Printf("Transaction signed and saved to: %s\n", outputFile)
		return nil
	},
}

var signMsgCmd = &cobra.Command{
	Use:   "message",
	Short: "Sign a message",
	Long:  `Sign an arbitrary message using a stored wallet key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load key
		manager, err := keystore.NewManager(keystoreDir)
		if err != nil {
			return fmt.Errorf("failed to create keystore manager: %v", err)
		}

		encryptedKey, err := manager.LoadKey(keyName)
		if err != nil {
			return fmt.Errorf("failed to load key: %v", err)
		}

		// Decrypt key
		privateKey, err := keystore.DecryptKey(encryptedKey, password)
		if err != nil {
			return fmt.Errorf("failed to decrypt key: %v", err)
		}

		// Sign message
		signature, err := core.SignMessage([]byte(message), privateKey)
		if err != nil {
			return fmt.Errorf("failed to sign message: %v", err)
		}

		// Write output
		if err := ioutil.WriteFile(outputFile, []byte(signature), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %v", err)
		}

		fmt.Printf("Message signed and saved to: %s\n", outputFile)
		return nil
	},
}

func init() {
	// Add flags
	SignCmd.PersistentFlags().StringVar(&keystoreDir, "keystore", ".keystore", "Keystore directory")
	SignCmd.PersistentFlags().StringVar(&keyName, "name", "", "Key name")
	SignCmd.PersistentFlags().StringVar(&password, "password", "", "Key password")
	SignCmd.PersistentFlags().StringVar(&outputFile, "output", "", "Output file")

	signTxCmd.Flags().StringVar(&inputFile, "input", "", "Input transaction file")
	signTxCmd.Flags().StringVar(&chainName, "chain", "ethereum", "Chain name")

	signMsgCmd.Flags().StringVar(&message, "message", "", "Message to sign")

	// Mark required flags
	SignCmd.MarkPersistentFlagRequired("name")
	SignCmd.MarkPersistentFlagRequired("password")
	SignCmd.MarkPersistentFlagRequired("output")

	signTxCmd.MarkFlagRequired("input")
	signMsgCmd.MarkFlagRequired("message")

	// Add commands
	SignCmd.AddCommand(signTxCmd)
	SignCmd.AddCommand(signMsgCmd)
}
