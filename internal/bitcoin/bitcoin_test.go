package bitcoin

import (
	"encoding/hex"
	"testing"

	"anvil/internal/crypto"
)

// Test vectors for Bitcoin address derivation
var bitcoinTestVectors = []struct {
	name     string
	mnemonic string
	path     string
	expected map[string]string // coin -> address
}{
	{
		name:     "Standard BIP39 test vector",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		path:     "m/44'/0'/0'/0/0",
		expected: map[string]string{
			"BTC": "1LqBGSKuX5yYUonjxT5qGfpUsXKYYWeabA",
		},
	},
	{
		name:     "Second address",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		path:     "m/44'/0'/0'/0/1",
		expected: map[string]string{
			"BTC": "1Ak8PffB2meyfYnbXZR9EGfLfFZVpzJvQP",
		},
	},
	{
		name:     "SegWit address",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		path:     "m/84'/0'/0'/0/0",
		expected: map[string]string{
			"BTC": "1JaUQDVNRdhfNsVncGkXedaPSM5Gc54Hso",
		},
	},
}

func TestBitcoinAddressGeneration(t *testing.T) {
	for _, tv := range bitcoinTestVectors {
		t.Run(tv.name, func(t *testing.T) {
			// Convert mnemonic to seed
			seed, err := crypto.MnemonicToSeed(tv.mnemonic, "")
			if err != nil {
				t.Fatalf("Failed to generate seed: %v", err)
			}
			defer crypto.SecureZeroMemory(seed)

			// Create Bitcoin coin instance
			btc := NewBitcoin()
			
			// Derive account
			account, err := btc.DeriveAccount(seed, tv.path)
			if err != nil {
				t.Fatalf("Failed to derive account: %v", err)
			}

			// Check expected address
			expectedAddr := tv.expected["BTC"]
			if account.Address != expectedAddr {
				t.Errorf("Address mismatch for %s:\nExpected: %s\nActual:   %s", 
					tv.path, expectedAddr, account.Address)
			}

			// Validate the generated address
			if !btc.ValidateAddress(account.Address) {
				t.Errorf("Generated address failed validation: %s", account.Address)
			}

			// Test path consistency
			if account.Path != tv.path {
				t.Errorf("Path mismatch: expected %s, got %s", tv.path, account.Path)
			}

			// Test symbol
			if account.Symbol != "BTC" {
				t.Errorf("Symbol mismatch: expected BTC, got %s", account.Symbol)
			}
		})
	}
}

func TestDogecoinAddressGeneration(t *testing.T) {
	// Test Dogecoin with known vectors
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, err := crypto.MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to generate seed: %v", err)
	}
	defer crypto.SecureZeroMemory(seed)

	doge := NewDogecoin()
	account, err := doge.DeriveAccount(seed, "m/44'/3'/0'/0/0")
	if err != nil {
		t.Fatalf("Failed to derive Dogecoin account: %v", err)
	}

	// Dogecoin addresses should start with 'D'
	if account.Address[0] != 'D' {
		t.Errorf("Dogecoin address should start with 'D', got: %s", account.Address)
	}

	// Should be valid
	if !doge.ValidateAddress(account.Address) {
		t.Errorf("Generated Dogecoin address failed validation: %s", account.Address)
	}

	// Test symbol
	if account.Symbol != "DOGE" {
		t.Errorf("Symbol mismatch: expected DOGE, got %s", account.Symbol)
	}
}

func TestBitcoinWIF(t *testing.T) {
	btc := NewBitcoin()
	
	// Test with known private key
	privateKeyHex := "0000000000000000000000000000000000000000000000000000000000000001"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	wif, err := btc.PrivateKeyToWIF(privateKeyBytes)
	if err != nil {
		t.Fatalf("Failed to generate WIF: %v", err)
	}

	// WIF should be reasonable length and format
	if len(wif) < 50 || len(wif) > 55 {
		t.Errorf("WIF length seems wrong: %d characters", len(wif))
	}

	// Should start with expected character for mainnet compressed WIF
	if wif[0] != 'K' && wif[0] != 'L' {
		t.Errorf("WIF should start with 'K' or 'L' for compressed mainnet, got: %c", wif[0])
	}
}

func TestBitcoinStandardPaths(t *testing.T) {
	btc := NewBitcoin()
	paths := btc.GetStandardDerivationPaths()
	
	expectedPaths := []string{
		"m/44'/0'/0'/0/0",  // BIP44 Legacy
		"m/49'/0'/0'/0/0",  // BIP49 P2SH-P2WPKH
		"m/84'/0'/0'/0/0",  // BIP84 Native SegWit
	}

	if len(paths) != len(expectedPaths) {
		t.Errorf("Expected %d standard paths, got %d", len(expectedPaths), len(paths))
	}

	for i, expectedPath := range expectedPaths {
		if i < len(paths) && paths[i] != expectedPath {
			t.Errorf("Path %d mismatch: expected %s, got %s", i, expectedPath, paths[i])
		}
	}
}

func TestCoinType(t *testing.T) {
	btc := NewBitcoin()
	if btc.GetCoinType() != 0 {
		t.Errorf("Bitcoin coin type should be 0, got %d", btc.GetCoinType())
	}

	doge := NewDogecoin()
	if doge.GetCoinType() != 3 {
		t.Errorf("Dogecoin coin type should be 3, got %d", doge.GetCoinType())
	}
}

// Benchmark tests for performance
func BenchmarkBitcoinAddressGeneration(b *testing.B) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := crypto.MnemonicToSeed(mnemonic, "")
	defer crypto.SecureZeroMemory(seed)
	
	btc := NewBitcoin()
	path := "m/44'/0'/0'/0/0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := btc.DeriveAccount(seed, path)
		if err != nil {
			b.Fatalf("Failed to derive account: %v", err)
		}
	}
}
