package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// Coin represents a cryptocurrency that can generate accounts
type Coin interface {
	Name() string
	Symbol() string
	DeriveAccount(seed []byte, path string) (Account, error)
}

// Account represents a cryptocurrency account with keys and address
type Account struct {
	Path       string    `json:"path"`
	PrivateKey []byte    `json:"private_key,omitempty"` // Omit in JSON unless requested
	PublicKey  []byte    `json:"public_key"`
	Address    string    `json:"address"`
	Symbol     string    `json:"symbol"`
	CreatedAt  time.Time `json:"created_at"`
}

// Wallet holds the master seed and derived accounts
type Wallet struct {
	Mnemonic  string              `json:"mnemonic,omitempty"` // Omit unless requested
	Seed      []byte              `json:"seed,omitempty"`     // Omit unless requested
	Accounts  []Account           `json:"accounts"`
	CoinTypes map[string][]uint32 `json:"coin_types"` // Maps symbol to BIP44 coin types
	CreatedAt time.Time           `json:"created_at"`
	Version   string              `json:"version"`
}

// MarshalSafeJSON returns JSON without sensitive fields
func (w *Wallet) MarshalSafeJSON() ([]byte, error) {
	safe := struct {
		Accounts  []SafeAccount       `json:"accounts"`
		CoinTypes map[string][]uint32 `json:"coin_types"`
		CreatedAt time.Time           `json:"created_at"`
		Version   string              `json:"version"`
	}{
		Accounts:  make([]SafeAccount, len(w.Accounts)),
		CoinTypes: w.CoinTypes,
		CreatedAt: w.CreatedAt,
		Version:   w.Version,
	}

	for i, acc := range w.Accounts {
		safe.Accounts[i] = SafeAccount{
			Path:      acc.Path,
			PublicKey: acc.PublicKey,
			Address:   acc.Address,
			Symbol:    acc.Symbol,
			CreatedAt: acc.CreatedAt,
		}
	}

	return json.Marshal(safe)
}

// SafeAccount represents an account without sensitive information
type SafeAccount struct {
	Path      string    `json:"path"`
	PublicKey []byte    `json:"public_key"`
	Address   string    `json:"address"`
	Symbol    string    `json:"symbol"`
	CreatedAt time.Time `json:"created_at"`
}

// DerivationPath represents a BIP32/BIP44 derivation path
type DerivationPath struct {
	Purpose  uint32 `json:"purpose"`   // Usually 44, 49, or 84
	CoinType uint32 `json:"coin_type"` // BIP44 registered coin type
	Account  uint32 `json:"account"`   // Account index
	Change   uint32 `json:"change"`    // 0 for external, 1 for internal
	Index    uint32 `json:"index"`     // Address index
}

// String returns the string representation of the derivation path
func (dp DerivationPath) String() string {
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d", dp.Purpose, dp.CoinType, dp.Account, dp.Change, dp.Index)
}

// ParseDerivationPath parses a string like "m/44'/0'/0'/0/0" into DerivationPath
func ParseDerivationPath(path string) (DerivationPath, error) {
	// Implementation will be in crypto package
	return DerivationPath{}, nil
}

// OutputFormat specifies how wallet data should be formatted
type OutputFormat int

const (
	OutputJSON OutputFormat = iota
	OutputText
	OutputQR
	OutputPaper
)

// OutputOptions controls what information is included in output
type OutputOptions struct {
	Format          OutputFormat `json:"format"`
	IncludePrivate  bool         `json:"include_private"`
	IncludeMnemonic bool         `json:"include_mnemonic"`
	IncludeQR       bool         `json:"include_qr"`
	FilePath        string       `json:"file_path,omitempty"`
}
