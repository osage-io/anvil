package solana

import (
	"crypto/ed25519"
	"testing"

	"anvil/internal/crypto"
	"github.com/mr-tron/base58"
)

func TestSolanaCoin_Name(t *testing.T) {
	coin := NewSolana()
	if coin.Name() != "Solana" {
		t.Errorf("Expected name 'Solana', got '%s'", coin.Name())
	}
}

func TestSolanaCoin_Symbol(t *testing.T) {
	coin := NewSolana()
	if coin.Symbol() != "SOL" {
		t.Errorf("Expected symbol 'SOL', got '%s'", coin.Symbol())
	}
}

func TestSolanaCoin_GetCoinType(t *testing.T) {
	coin := NewSolana()
	if coin.GetCoinType() != 501 {
		t.Errorf("Expected coin type 501, got %d", coin.GetCoinType())
	}
}

func TestSolanaCoin_ValidateAddress(t *testing.T) {
	coin := NewSolana()

	// Test valid address
	validAddress := "11111111111111111111111111111112" // System program ID
	if !coin.ValidateAddress(validAddress) {
		t.Errorf("Valid address %s was marked as invalid", validAddress)
	}

	// Test invalid address (too short)
	invalidAddress := "invalid"
	if coin.ValidateAddress(invalidAddress) {
		t.Errorf("Invalid address %s was marked as valid", invalidAddress)
	}
}

func TestSolanaCoin_GetStandardDerivationPaths(t *testing.T) {
	coin := NewSolana()
	paths := coin.GetStandardDerivationPaths()

	expected := []string{
		"m/44'/501'/0'/0'",
		"m/44'/501'/1'/0'",
		"m/44'/501'/2'/0'",
	}

	if len(paths) != len(expected) {
		t.Errorf("Expected %d paths, got %d", len(expected), len(paths))
	}

	for i, path := range paths {
		if path != expected[i] {
			t.Errorf("Expected path %s at index %d, got %s", expected[i], i, path)
		}
	}
}

func TestSolanaCoin_DeriveAccount(t *testing.T) {
	coin := NewSolana()

	// Generate a test seed
	mnemonic, err := crypto.GenerateMnemonic(128)
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	seed, err := crypto.MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to convert mnemonic to seed: %v", err)
	}
	defer crypto.ClearBytes(seed)

	// Derive account
	path := "m/44'/501'/0'/0'"
	account, err := coin.DeriveAccount(seed, path)
	if err != nil {
		t.Fatalf("Failed to derive account: %v", err)
	}

	// Validate account fields
	if account.Symbol != "SOL" {
		t.Errorf("Expected symbol 'SOL', got '%s'", account.Symbol)
	}

	if account.Path != path {
		t.Errorf("Expected path '%s', got '%s'", path, account.Path)
	}

	if len(account.PublicKey) != ed25519.PublicKeySize {
		t.Errorf("Expected public key size %d, got %d", ed25519.PublicKeySize, len(account.PublicKey))
	}

	if len(account.PrivateKey) != ed25519.PrivateKeySize {
		t.Errorf("Expected private key seed size %d, got %d", ed25519.PrivateKeySize, len(account.PrivateKey))
	}

	// Validate that address is base58 encoded public key
	expectedAddress := base58.Encode(account.PublicKey)
	if account.Address != expectedAddress {
		t.Errorf("Expected address '%s', got '%s'", expectedAddress, account.Address)
	}

	// Validate address format
	if !coin.ValidateAddress(account.Address) {
		t.Errorf("Generated address %s is not valid", account.Address)
	}
}

func TestSolanaNetworkVariants(t *testing.T) {
	mainnet := NewSolana()
	testnet := NewSolanaTestnet()
	devnet := NewSolanaDevnet()

	if mainnet.GetNetwork() != "mainnet-beta" {
		t.Errorf("Expected mainnet network 'mainnet-beta', got '%s'", mainnet.GetNetwork())
	}

	if testnet.GetNetwork() != "testnet" {
		t.Errorf("Expected testnet network 'testnet', got '%s'", testnet.GetNetwork())
	}

	if devnet.GetNetwork() != "devnet" {
		t.Errorf("Expected devnet network 'devnet', got '%s'", devnet.GetNetwork())
	}

	// All should have same coin type and symbol
	if mainnet.GetCoinType() != testnet.GetCoinType() || mainnet.GetCoinType() != devnet.GetCoinType() {
		t.Error("All Solana networks should have the same coin type")
	}

	if mainnet.Symbol() != testnet.Symbol() || mainnet.Symbol() != devnet.Symbol() {
		t.Error("All Solana networks should have the same symbol")
	}
}
