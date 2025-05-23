package main

import (
	"fmt"
	"os"

	"github.com/aryehky/gosignervaultcli/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gosignervaultcli",
	Short: "GoSignerVaultCLI - A secure CLI wallet and transaction signer",
	Long: `GoSignerVaultCLI is a lightweight, secure, and extensible command-line interface (CLI) wallet
and transaction signer built in Go. It allows you to securely generate and manage private keys
offline, sign transactions for Ethereum-compatible blockchains, and export signed payloads for broadcast.`,
}

func init() {
	// Add commands
	rootCmd.AddCommand(cmd.KeysCmd)
	rootCmd.AddCommand(cmd.SignCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
