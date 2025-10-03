package tron

import (
	"strings"
	"testing"

	"anvil/internal/crypto"
)

// Test vectors for TRON address derivation
var tronTestVectors = []struct {
	name     string
	mnemonic string
	path     string
	expected string
}{
	{
		name:     "Standard BIP39 test vector",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		path:     "m/44'/195'/0'/0/0",
		expected: "TUEZSdKsoDHQMeZwihtdoBiN46zxhGWYdH",
	},
}

func TestTronAddressGeneration(t *testing.T) {
	for _, tv := range tronTestVectors {
		t.Run(tv.name, func(t *testing.T) {
			// Convert mnemonic to seed
			seed, err := crypto.MnemonicToSeed(tv.mnemonic, "")
			if err != nil {
				t.Fatalf("Failed to generate seed: %v", err)
			}
			defer crypto.SecureZeroMemory(seed)

			// Create TRON coin instance
			trx := NewTron()
			
			// Derive account
			account, err := trx.DeriveAccount(seed, tv.path)
			if err != nil {
				t.Fatalf("Failed to derive account: %v", err)
			}

			// Check expected address
			if account.Address != tv.expected {
				t.Errorf("Address mismatch for %s:\nExpected: %s\nActual:   %s", 
					tv.path, tv.expected, account.Address)
			}

			// Validate the generated address
			if !trx.ValidateAddress(account.Address) {
				t.Errorf("Generated address failed validation: %s", account.Address)
			}

			// TRON addresses should start with 'T'
			if !strings.HasPrefix(account.Address, "T") {
				t.Errorf("TRON address should start with 'T', got: %s", account.Address)
			}

			// Test path consistency
			if account.Path != tv.path {
				t.Errorf("Path mismatch: expected %s, got %s", tv.path, account.Path)
			}

			// Test symbol
			if account.Symbol != "TRX" {
				t.Errorf("Symbol mismatch: expected TRX, got %s", account.Symbol)
			}
		})
	}
}

func TestTronAddressValidation(t *testing.T) {
	trx := NewTron()
	
	validAddresses := []string{
		"TUEZSdKsoDHQMeZwihtdoBiN46zxhGWYdH",
		"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", // USDT on TRON
		"TLyqzVGLV1srkB7dToTAEqgDSfPtXRJZYH",
	}
	
	for _, addr := range validAddresses {
		if !trx.ValidateAddress(addr) {
			t.Errorf("Valid TRON address rejected: %s", addr)
		}
	}
	
	invalidAddresses := []string{
		"1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2", // Bitcoin address
		"0x742D35CC6634C0532925A3b8D4b72866", // Ethereum address
		"TUEZSdKsoDHQMeZwihtdoBiN46zxhGWYd",  // Too short
		"TUEZSdKsoDHQMeZwihtdoBiN46zxhGWYdHH", // Too long
		"AUEZSdKsoDHQMeZwihtdoBiN46zxhGWYdH", // Wrong prefix
		"", // Empty
	}
	
	for _, addr := range invalidAddresses {
		if trx.ValidateAddress(addr) {
			t.Errorf("Invalid TRON address accepted: %s", addr)
		}
	}
}

func TestTronAddressConversion(t *testing.T) {
	trx := NewTron()
	
	testCases := []struct {
		tronAddress string
		hexAddress  string
	}{
		{
			tronAddress: "TPhbUUoHYFjf81yex3DuRwvvBR2ZSQmLk1",
			hexAddress:  "0x969ddd6b04052f60be05c9ee7ae228dafec5c9e5",
		},
	}
	
	for _, tc := range testCases {
		// Test TRON to hex conversion
		hexResult, err := trx.AddressToHex(tc.tronAddress)
		if err != nil {
			t.Errorf("Failed to convert TRON address to hex: %v", err)
		}
		
		if !strings.EqualFold(hexResult, tc.hexAddress) {
			t.Errorf("Hex conversion mismatch:\nTRON: %s\nExpected: %s\nActual:   %s", 
				tc.tronAddress, tc.hexAddress, hexResult)
		}
		
		// Test hex to TRON conversion
		tronResult, err := trx.HexToAddress(tc.hexAddress)
		if err != nil {
			t.Errorf("Failed to convert hex address to TRON: %v", err)
		}
		
		if tronResult != tc.tronAddress {
			t.Errorf("TRON conversion mismatch:\nHex: %s\nExpected: %s\nActual:   %s", 
				tc.hexAddress, tc.tronAddress, tronResult)
		}
	}
}

func TestTronStandardPaths(t *testing.T) {
	trx := NewTron()
	paths := trx.GetStandardDerivationPaths()
	
	expectedPaths := []string{
		"m/44'/195'/0'/0/0", // BIP44 standard path
		"m/44'/195'/0'/0/1", // Second address
		"m/44'/195'/1'/0/0", // Change addresses
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

func TestTronCoinType(t *testing.T) {
	trx := NewTron()
	if trx.GetCoinType() != 195 {
		t.Errorf("TRON coin type should be 195, got %d", trx.GetCoinType())
	}
}

func TestDoubleSHA256(t *testing.T) {
	// Test the double SHA256 function used for TRON checksums
	testData := []byte("test")
	result := doubleSHA256(testData)
	
	if len(result) != 32 {
		t.Errorf("Double SHA256 should return 32 bytes, got %d", len(result))
	}
	
	// Test with known input
	knownInput := []byte{0x41, 0x01, 0x02, 0x03, 0x04, 0x05}
	result1 := doubleSHA256(knownInput)
	result2 := doubleSHA256(knownInput)
	
	// Should be deterministic
	if len(result1) != len(result2) {
		t.Error("Double SHA256 results have different lengths")
	}
	
	for i, b := range result1 {
		if b != result2[i] {
			t.Error("Double SHA256 is not deterministic")
			break
		}
	}
}

// Benchmark tests for performance
func BenchmarkTronAddressGeneration(b *testing.B) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := crypto.MnemonicToSeed(mnemonic, "")
	defer crypto.SecureZeroMemory(seed)
	
	trx := NewTron()
	path := "m/44'/195'/0'/0/0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := trx.DeriveAccount(seed, path)
		if err != nil {
			b.Fatalf("Failed to derive account: %v", err)
		}
	}
}

func BenchmarkTronAddressValidation(b *testing.B) {
	trx := NewTron()
	address := "TUEZSdKsoDHQMeZwihtdoBiN46zxhGWYdH"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trx.ValidateAddress(address)
	}
}

func BenchmarkDoubleSHA256(b *testing.B) {
	data := make([]byte, 21)
	data[0] = 0x41 // TRON prefix
	for i := 1; i < 21; i++ {
		data[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doubleSHA256(data)
	}
}
