package ethereum

import (
	"strings"
	"testing"

	"anvil/internal/crypto"
)

// Test vectors for Ethereum address derivation
var ethereumTestVectors = []struct {
	name     string
	mnemonic string
	path     string
	expected string
}{
	{
		name:     "Standard BIP39 test vector",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		path:     "m/44'/60'/0'/0/0",
		expected: "0x9858EfFD232B4033E47d90003D41EC34EcaEda94",
	},
	{
		name:     "Second address",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		path:     "m/44'/60'/0'/0/1",
		expected: "0x6Fac4D18c912343BF86fa7049364Dd4E424Ab9C0",
	},
}

func TestEthereumAddressGeneration(t *testing.T) {
	for _, tv := range ethereumTestVectors {
		t.Run(tv.name, func(t *testing.T) {
			// Convert mnemonic to seed
			seed, err := crypto.MnemonicToSeed(tv.mnemonic, "")
			if err != nil {
				t.Fatalf("Failed to generate seed: %v", err)
			}
			defer crypto.SecureZeroMemory(seed)

			// Create Ethereum coin instance
			eth := NewEthereum()
			
			// Derive account
			account, err := eth.DeriveAccount(seed, tv.path)
			if err != nil {
				t.Fatalf("Failed to derive account: %v", err)
			}

			// Check expected address (case-insensitive comparison)
			if !strings.EqualFold(account.Address, tv.expected) {
				t.Errorf("Address mismatch for %s:\nExpected: %s\nActual:   %s", 
					tv.path, tv.expected, account.Address)
			}

			// Validate the generated address
			if !eth.ValidateAddress(account.Address) {
				t.Errorf("Generated address failed validation: %s", account.Address)
			}

			// Test EIP-55 checksum (should be properly formatted)
			if account.Address != eth.toChecksumAddress(strings.ToLower(account.Address)) {
				t.Errorf("Address is not properly checksummed: %s", account.Address)
			}

			// Test path consistency
			if account.Path != tv.path {
				t.Errorf("Path mismatch: expected %s, got %s", tv.path, account.Path)
			}

			// Test symbol
			if account.Symbol != "ETH" {
				t.Errorf("Symbol mismatch: expected ETH, got %s", account.Symbol)
			}
		})
	}
}

func TestBinanceCoinAddressGeneration(t *testing.T) {
	// Test BNB Smart Chain with known vectors
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, err := crypto.MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to generate seed: %v", err)
	}
	defer crypto.SecureZeroMemory(seed)

	bnb := NewBinanceCoin()
	account, err := bnb.DeriveAccount(seed, "m/44'/60'/0'/0/0")
	if err != nil {
		t.Fatalf("Failed to derive BNB account: %v", err)
	}

	// BNB addresses should start with '0x'
	if !strings.HasPrefix(account.Address, "0x") {
		t.Errorf("BNB address should start with '0x', got: %s", account.Address)
	}

	// Should be valid
	if !bnb.ValidateAddress(account.Address) {
		t.Errorf("Generated BNB address failed validation: %s", account.Address)
	}

	// Test symbol
	if account.Symbol != "BNB" {
		t.Errorf("Symbol mismatch: expected BNB, got %s", account.Symbol)
	}

	// Should be same as Ethereum address for same derivation path
	eth := NewEthereum()
	ethAccount, err := eth.DeriveAccount(seed, "m/44'/60'/0'/0/0")
	if err != nil {
		t.Fatalf("Failed to derive ETH account: %v", err)
	}
	
	if !strings.EqualFold(account.Address, ethAccount.Address) {
		t.Errorf("BNB and ETH addresses should be the same for same path:\nBNB: %s\nETH: %s", 
			account.Address, ethAccount.Address)
	}
}

func TestEthereumAddressValidation(t *testing.T) {
	eth := NewEthereum()
	
	validAddresses := []string{
		"0x9858EfFD232B4033E47d90003D41EC34EcaEda94", // Checksummed
		"0x9858effd232b4033e47d90003d41ec34ecaeda94", // Lowercase
		"0x9858EFFD232B4033E47D90003D41EC34ECAEDA94", // Uppercase
	}
	
	for _, addr := range validAddresses {
		if !eth.ValidateAddress(addr) {
			t.Errorf("Valid address rejected: %s", addr)
		}
	}
	
	invalidAddresses := []string{
		"9858effd232b4033e47d90003d41ec34ecaeda94", // Missing 0x
		"0x9858effd232b4033e47d90003d41ec34ecaeda9",  // Too short
		"0x9858effd232b4033e47d90003d41ec34ecaeda944", // Too long
		"0x9858effd232b4033e47d90003d41ec34ecaeda9g", // Invalid hex
		"0x9858EfFd232b4033E47d90003D41EC34eCAeDA95", // Wrong checksum
	}
	
	for _, addr := range invalidAddresses {
		if eth.ValidateAddress(addr) {
			t.Errorf("Invalid address accepted: %s", addr)
		}
	}
}

func TestEIP55Checksum(t *testing.T) {
	eth := NewEthereum()
	
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "0x9858effd232b4033e47d90003d41ec34ecaeda94",
			expected: "0x9858EfFD232B4033E47d90003D41EC34EcaEda94",
		},
		{
			input:    "0x6fac4d18c912343bf86fa7049364dd4e424ab9c0",
			expected: "0x6Fac4D18c912343BF86fa7049364Dd4E424Ab9C0",
		},
	}
	
	for _, tc := range testCases {
		result := eth.toChecksumAddress(tc.input)
		if result != tc.expected {
			t.Errorf("Checksum mismatch:\nInput:    %s\nExpected: %s\nActual:   %s", 
				tc.input, tc.expected, result)
		}
	}
}

func TestEthereumStandardPaths(t *testing.T) {
	eth := NewEthereum()
	paths := eth.GetStandardDerivationPaths()
	
	expectedPaths := []string{
		"m/44'/60'/0'/0/0",  // BIP44 standard path
		"m/44'/60'/0'/0/1",  // Second address
		"m/44'/60'/1'/0/0",  // Change addresses
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

func TestCoinTypes(t *testing.T) {
	eth := NewEthereum()
	if eth.GetCoinType() != 60 {
		t.Errorf("Ethereum coin type should be 60, got %d", eth.GetCoinType())
	}

	bnb := NewBinanceCoin()
	if bnb.GetCoinType() != 60 {
		t.Errorf("BNB coin type should be 60 (same as Ethereum), got %d", bnb.GetCoinType())
	}
	
	// Test chain IDs
	if eth.GetChainID() != 1 {
		t.Errorf("Ethereum chain ID should be 1, got %d", eth.GetChainID())
	}
	
	if bnb.GetChainID() != 56 {
		t.Errorf("BNB Smart Chain ID should be 56, got %d", bnb.GetChainID())
	}
}

// Benchmark tests for performance
func BenchmarkEthereumAddressGeneration(b *testing.B) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := crypto.MnemonicToSeed(mnemonic, "")
	defer crypto.SecureZeroMemory(seed)
	
	eth := NewEthereum()
	path := "m/44'/60'/0'/0/0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := eth.DeriveAccount(seed, path)
		if err != nil {
			b.Fatalf("Failed to derive account: %v", err)
		}
	}
}

func BenchmarkEIP55Checksum(b *testing.B) {
	eth := NewEthereum()
	address := "0x9858effd232b4033e47d90003d41ec34ecaeda94"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eth.toChecksumAddress(address)
	}
}
