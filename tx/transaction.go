package tx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Transaction represents an Ethereum transaction
type Transaction struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Value    *big.Int        `json:"value"`
	Gas      uint64          `json:"gas"`
	GasPrice *big.Int        `json:"gasPrice"`
	Data     []byte          `json:"data"`
	Nonce    uint64          `json:"nonce"`
	ChainID  *big.Int        `json:"chainId"`
}

// ToEthereumTx converts the Transaction to an Ethereum types.Transaction
func (t *Transaction) ToEthereumTx() *types.Transaction {
	return types.NewTransaction(
		t.Nonce,
		*t.To,
		t.Value,
		t.Gas,
		t.GasPrice,
		t.Data,
	)
}

// FromEthereumTx creates a Transaction from an Ethereum types.Transaction
func FromEthereumTx(tx *types.Transaction, from common.Address) *Transaction {
	return &Transaction{
		From:     from,
		To:       tx.To(),
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Data:     tx.Data(),
		Nonce:    tx.Nonce(),
		ChainID:  tx.ChainId(),
	}
}

// ToRLP encodes the transaction to RLP format
func (t *Transaction) ToRLP() ([]byte, error) {
	ethTx := t.ToEthereumTx()
	return ethTx.MarshalBinary()
}

// FromRLP decodes an RLP-encoded transaction
func FromRLP(data []byte) (*Transaction, error) {
	var tx types.Transaction
	if err := tx.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return FromEthereumTx(&tx, common.Address{}), nil
}
