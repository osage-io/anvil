package solana

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"anvil/internal/crypto"
	"anvil/pkg/types"
	"github.com/blocto/solana-go-sdk/pkg/hdwallet"
	solanatypes "github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
)

// SolanaCoin implements the types.Coin interface for Solana
type SolanaCoin struct {
	name     string
	symbol   string
	coinType uint32
	network  string
}

// NewSolana creates a new Solana coin instance
func NewSolana() *SolanaCoin {
	return &SolanaCoin{
		name:     "Solana",
		symbol:   "SOL",
		coinType: 501, // BIP44 coin type for Solana
		network:  "mainnet-beta",
	}
}

// NewSolanaTestnet creates a new Solana testnet instance
func NewSolanaTestnet() *SolanaCoin {
	return &SolanaCoin{
		name:     "Solana Testnet",
		symbol:   "SOL",
		coinType: 501,
		network:  "testnet",
	}
}

// NewSolanaDevnet creates a new Solana devnet instance
func NewSolanaDevnet() *SolanaCoin {
	return &SolanaCoin{
		name:     "Solana Devnet",
		symbol:   "SOL",
		coinType: 501,
		network:  "devnet",
	}
}

// Name returns the full name of the cryptocurrency
func (s *SolanaCoin) Name() string {
	return s.name
}

// Symbol returns the symbol/ticker of the cryptocurrency
func (s *SolanaCoin) Symbol() string {
	return s.symbol
}

// DeriveAccount derives a new account for the given seed and derivation path
func (s *SolanaCoin) DeriveAccount(seed []byte, path string) (types.Account, error) {
	// Use Solana's ED25519 HD wallet derivation
	// Solana uses paths like m/44'/501'/0'/0' where all components are hardened
	key, err := hdwallet.Derived(path, seed)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to derive Solana key: %w", err)
	}

	// Create Solana account from the derived private key
	account, err := solanatypes.AccountFromSeed(key.PrivateKey)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to create Solana account: %w", err)
	}

	// Get the address as base58 string
	address := account.PublicKey.ToBase58()

	solanaAccount := types.Account{
		Path:       path,
		PrivateKey: account.PrivateKey, // Full 64-byte private key
		PublicKey:  account.PublicKey.Bytes(),
		Address:    address,
		Symbol:     s.symbol,
		CreatedAt:  time.Now(),
	}

	// Clear sensitive key data from memory
	crypto.SecureZeroMemory(key.PrivateKey)

	return solanaAccount, nil
}

// ValidateAddress checks if an address is valid for Solana
func (s *SolanaCoin) ValidateAddress(address string) bool {
	// Try to decode as base58 and check length
	decoded, err := base58.Decode(address)
	if err != nil {
		return false
	}
	// Solana public keys are exactly 32 bytes
	return len(decoded) == 32
}

// GetStandardDerivationPaths returns common derivation paths for Solana
func (s *SolanaCoin) GetStandardDerivationPaths() []string {
	coinType := s.coinType
	return []string{
		fmt.Sprintf("m/44'/%d'/0'/0'", coinType), // Standard Solana path (hardened)
		fmt.Sprintf("m/44'/%d'/1'/0'", coinType), // Second account
		fmt.Sprintf("m/44'/%d'/2'/0'", coinType), // Third account
	}
}

// GetCoinType returns the BIP44 coin type for Solana
func (s *SolanaCoin) GetCoinType() uint32 {
	return s.coinType
}

// GetNetwork returns the network name (mainnet-beta, testnet, devnet)
func (s *SolanaCoin) GetNetwork() string {
	return s.network
}

// PrivateKeyToBytes converts a Solana private key to bytes
func (s *SolanaCoin) PrivateKeyToBytes(privateKey ed25519.PrivateKey) []byte {
	return privateKey
}

// PublicKeyFromPrivate derives the public key from a private key
func (s *SolanaCoin) PublicKeyFromPrivate(privateKey ed25519.PrivateKey) ed25519.PublicKey {
	return privateKey.Public().(ed25519.PublicKey)
}

// AddressFromPublicKey converts a public key to a Solana address
func (s *SolanaCoin) AddressFromPublicKey(publicKey ed25519.PublicKey) string {
	return base58.Encode(publicKey)
}
