# 🔐 GoSignerVaultCLI: Go-Powered Cold Wallet CLI & Transaction Signer

**GoSignerVaultCLI** is a lightweight, secure, and extensible command-line interface (CLI) wallet and transaction signer built in **Go**. Designed for crypto enthusiasts, developers, and validators, GoSignerVaultCLI allows you to securely generate and manage private keys offline, sign transactions for Ethereum-compatible blockchains, and export signed payloads for broadcast.

Ideal for air-gapped environments or minimal setups, GoSignerVaultCLI focuses on **security-first architecture**, **zero external dependencies**, and **extensibility for any EVM-based chain**.

---

## ✨ Features

* 🔑 **Generate Cold Wallet Keys**
  Securely generate new wallets (private/public keypairs) using Go's `crypto/ecdsa` and save them to encrypted keystore files.

* 🛡️ **Offline Transaction Signing**
  Import or paste unsigned transactions, sign them locally, and export raw signed transactions to be broadcast separately.

* 🔗 **Ethereum & EVM-Compatible**
  Full support for Ethereum, Polygon, BNB Smart Chain, Avalanche C-Chain, etc. via customizable chain configs.

* 📁 **Keystore Encryption**
  Encrypt private keys using AES-256 and store them locally in password-protected JSON files.

* 🧩 **Modular Chain Configs**
  Easily switch between supported networks or add your own by editing a simple TOML config.

* 🔋 **Message Signing (EIP-191)**
  Sign arbitrary messages using the `eth_sign` method for use in DApps, DAOs, and smart contract authentication.

---

## 📂 Project Structure

```
gosignervaultcli/
│
├── cmd/                   # CLI logic (Cobra-based)
│   ├── keys.go            # Commands to generate/list/show keys
│   ├── sign.go            # Commands to sign transactions/messages
│   └── utils.go           # Helper functions for CLI
│
├── core/                  # Core crypto logic
│   ├── wallet.go          # Key generation, encryption, decryption
│   ├── signer.go          # Transaction and message signing
│   └── chain_config.go    # EVM chain configuration struct
│
├── keystore/              # Encrypted key file manager
│   ├── encrypt.go         # AES-based encryption helpers
│   └── keystore.go        # Load/save encrypted keys
│
├── tx/                    # Transaction encoder/decoder
│   ├── txbuilder.go       # Raw transaction struct
│   └── broadcast.go       # Optional: broadcasting helper (for online setups)
│
├── go.mod
└── README.md
```

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/aryehky/gosignervaultcli
cd gosignervaultcli
go mod tidy
```

### 2. Build the CLI

```bash
go build -o gosignervaultcli main.go
```

### 3. Create a New Wallet

```bash
./gosignervaultcli keys generate --name mywallet
```

Enter a strong password to encrypt your keystore file.

### 4. Sign a Transaction (Offline)

```bash
./gosignervaultcli sign tx --input rawTx.json --wallet mywallet --output signedTx.json
```

### 5. Export for Broadcast

Upload the `signedTx.json` to an online machine and broadcast it with tools like [Etherscan Gas Tracker](https://etherscan.io/pushTx) or custom RPC broadcaster.

---

## 🛠 Configuration

Edit `config.toml` to define custom EVM chains:

```toml
[ethereum]
rpc_url = "https://mainnet.infura.io/v3/YOUR_INFURA_KEY"
chain_id = 1
symbol = "ETH"

[polygon]
rpc_url = "https://polygon-rpc.com"
chain_id = 137
symbol = "MATIC"
```

---

## 🧪 Test Coverage

Run unit tests for core modules:

```bash
go test ./core/...
go test ./keystore/...
```

---

## 🔒 Security Considerations

* **Never share your keystore files or passwords.**
* Keep this app on an **air-gapped device** for maximum cold storage protection.
* All private key handling is performed **in-memory** and securely zeroed after use.

---

## 📄 License

MIT License

---

## 🧠 Inspiration

GoSignerVaultCLI is inspired by cold wallet tools like Gnosis Safe CLI, Ledger Live, and MyCrypto desktop but reimagined in Go for transparency, auditability, and minimal dependency footprint.

---

## 🤝 Contributing

PRs welcome! Submit issues or feature requests under [Issues](https://github.com/aryehky/gosignervaultcli/issues).
