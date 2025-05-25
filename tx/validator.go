package tx

import (
	"fmt"
	"math/big"
)

// ValidationError represents a transaction validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validator handles transaction validation
type Validator struct {
	MinGasPrice *big.Int
	MaxGasPrice *big.Int
	MaxGasLimit uint64
	MinValue    *big.Int
	MaxValue    *big.Int
}

// NewValidator creates a new transaction validator
func NewValidator() *Validator {
	return &Validator{
		MinGasPrice: big.NewInt(1),             // 1 wei
		MaxGasPrice: big.NewInt(1000000000000), // 1000 gwei
		MaxGasLimit: 10000000,                  // 10M gas
		MinValue:    big.NewInt(0),
		MaxValue:    big.NewInt(0).Mul(big.NewInt(1000000), big.NewInt(1e18)), // 1M ETH
	}
}

// ValidateTransaction validates a transaction
func (v *Validator) ValidateTransaction(tx *Transaction) []ValidationError {
	var errors []ValidationError

	// Validate gas price
	if tx.GasPrice.Cmp(v.MinGasPrice) < 0 {
		errors = append(errors, ValidationError{
			Field:   "gasPrice",
			Message: fmt.Sprintf("gas price too low: %s < %s", tx.GasPrice.String(), v.MinGasPrice.String()),
		})
	}
	if tx.GasPrice.Cmp(v.MaxGasPrice) > 0 {
		errors = append(errors, ValidationError{
			Field:   "gasPrice",
			Message: fmt.Sprintf("gas price too high: %s > %s", tx.GasPrice.String(), v.MaxGasPrice.String()),
		})
	}

	// Validate gas limit
	if tx.Gas > v.MaxGasLimit {
		errors = append(errors, ValidationError{
			Field:   "gas",
			Message: fmt.Sprintf("gas limit too high: %d > %d", tx.Gas, v.MaxGasLimit),
		})
	}

	// Validate value
	if tx.Value.Cmp(v.MinValue) < 0 {
		errors = append(errors, ValidationError{
			Field:   "value",
			Message: fmt.Sprintf("value too low: %s < %s", tx.Value.String(), v.MinValue.String()),
		})
	}
	if tx.Value.Cmp(v.MaxValue) > 0 {
		errors = append(errors, ValidationError{
			Field:   "value",
			Message: fmt.Sprintf("value too high: %s > %s", tx.Value.String(), v.MaxValue.String()),
		})
	}

	// Validate address
	if tx.To == nil {
		errors = append(errors, ValidationError{
			Field:   "to",
			Message: "recipient address is required",
		})
	}

	// Validate chain ID
	if tx.ChainID == nil || tx.ChainID.Sign() <= 0 {
		errors = append(errors, ValidationError{
			Field:   "chainId",
			Message: "valid chain ID is required",
		})
	}

	return errors
}

// SetGasPriceLimits sets the minimum and maximum gas price
func (v *Validator) SetGasPriceLimits(min, max *big.Int) {
	v.MinGasPrice = min
	v.MaxGasPrice = max
}

// SetGasLimit sets the maximum gas limit
func (v *Validator) SetGasLimit(max uint64) {
	v.MaxGasLimit = max
}

// SetValueLimits sets the minimum and maximum transaction value
func (v *Validator) SetValueLimits(min, max *big.Int) {
	v.MinValue = min
	v.MaxValue = max
}
