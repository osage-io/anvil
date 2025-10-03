package crypto

import (
	"fmt"
	"strings"
	"testing"

	"anvil/pkg/types"
)

// Test vectors for BIP39 mnemonic generation and validation
func TestMnemonicGeneration(t *testing.T) {
	testCases := []struct {
		entropyBits int
		wordCount   int
	}{
		{128, 12},
		{160, 15},
		{192, 18},
		{224, 21},
		{256, 24},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_bits", tc.entropyBits), func(t *testing.T) {
			mnemonic, err := GenerateMnemonic(tc.entropyBits)
			if err != nil {
				t.Fatalf("Failed to generate mnemonic: %v", err)
			}

			// Check word count
			words := strings.Fields(mnemonic)
			if len(words) != tc.wordCount {
				t.Errorf("Expected %d words, got %d", tc.wordCount, len(words))
			}

			// Check if mnemonic is valid
			seed, err := MnemonicToSeed(mnemonic, "")
			if err != nil {
				t.Errorf("Generated mnemonic is invalid: %v", err)
			}

			// Seed should be 64 bytes (512 bits)
			if len(seed) != 64 {
				t.Errorf("Seed should be 64 bytes, got %d", len(seed))
			}

			// Clear the seed
			SecureZeroMemory(seed)
		})
	}
}

func TestMnemonicValidation(t *testing.T) {
	validMnemonics := []string{
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		"legal winner thank year wave sausage worth useful legal winner thank yellow",
		"letter advice cage absurd amount doctor acoustic avoid letter advice cage above",
	}

	for _, mnemonic := range validMnemonics {
		seed, err := MnemonicToSeed(mnemonic, "")
		if err != nil {
			t.Errorf("Valid mnemonic rejected: %s, error: %v", mnemonic, err)
		} else {
			SecureZeroMemory(seed)
		}
	}

	invalidMnemonics := []string{
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon",         // Too short
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon invalid", // Invalid word
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon", // Invalid checksum
		"",     // Empty
		"word", // Single word
	}

	for _, mnemonic := range invalidMnemonics {
		_, err := MnemonicToSeed(mnemonic, "")
		if err == nil {
			t.Errorf("Invalid mnemonic accepted: %s", mnemonic)
		}
	}
}

func TestMnemonicWithPassphrase(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	// Same mnemonic with different passphrases should produce different seeds
	seed1, err := MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to generate seed: %v", err)
	}
	defer SecureZeroMemory(seed1)

	seed2, err := MnemonicToSeed(mnemonic, "passphrase")
	if err != nil {
		t.Fatalf("Failed to generate seed with passphrase: %v", err)
	}
	defer SecureZeroMemory(seed2)

	// Seeds should be different
	if len(seed1) != len(seed2) {
		t.Error("Seeds have different lengths")
	} else {
		same := true
		for i, b := range seed1 {
			if b != seed2[i] {
				same = false
				break
			}
		}
		if same {
			t.Error("Seeds should be different with different passphrases")
		}
	}
}

func TestDerivationPathParsing(t *testing.T) {
	testCases := []struct {
		path     string
		expected types.DerivationPath
		valid    bool
	}{
		{
			path: "m/44'/0'/0'/0/0",
			expected: types.DerivationPath{
				Purpose:  44,
				CoinType: 0,
				Account:  0,
				Change:   0,
				Index:    0,
			},
			valid: true,
		},
		{
			path: "m/84'/0'/0'/0/5",
			expected: types.DerivationPath{
				Purpose:  84,
				CoinType: 0,
				Account:  0,
				Change:   0,
				Index:    5,
			},
			valid: true,
		},
		{
			path: "m/44'/60'/0'/0/1",
			expected: types.DerivationPath{
				Purpose:  44,
				CoinType: 60,
				Account:  0,
				Change:   0,
				Index:    1,
			},
			valid: true,
		},
		{
			path:  "44'/0'/0'/0/0", // Missing m/
			valid: false,
		},
		{
			path:  "m/44/0'/0'/0/0", // Purpose not hardened
			valid: false,
		},
		{
			path:  "m/44'/0/0'/0/0", // Coin type not hardened
			valid: false,
		},
		{
			path:  "m/44'/0'/0/0/0", // Account not hardened
			valid: false,
		},
		{
			path:  "m/44'/0'/0'/0'/0", // Index hardened (should not be)
			valid: false,
		},
		{
			path:  "m/44'/0'/0'/0", // Too few components
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result, err := ParseDerivationPath(tc.path)

			if tc.valid {
				if err != nil {
					t.Errorf("Valid path rejected: %s, error: %v", tc.path, err)
				} else {
					if result != tc.expected {
						t.Errorf("Path parsing mismatch:\nExpected: %+v\nActual:   %+v", tc.expected, result)
					}

					// Test string representation
					stringResult := result.String()
					if stringResult != tc.path {
						t.Errorf("String representation mismatch:\nExpected: %s\nActual:   %s", tc.path, stringResult)
					}
				}
			} else {
				if err == nil {
					t.Errorf("Invalid path accepted: %s", tc.path)
				}
			}
		})
	}
}

func TestKeyDerivation(t *testing.T) {
	// Test with known mnemonic and path
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, err := MnemonicToSeed(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to generate seed: %v", err)
	}
	defer SecureZeroMemory(seed)

	// Test key derivation
	key, err := DeriveKey(seed, "m/44'/0'/0'/0/0")
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	// Key should be 32 bytes
	if len(key.Key) != 32 {
		t.Errorf("Derived key should be 32 bytes, got %d", len(key.Key))
	}

	// Test that different paths produce different keys
	key2, err := DeriveKey(seed, "m/44'/0'/0'/0/1")
	if err != nil {
		t.Fatalf("Failed to derive second key: %v", err)
	}

	// Keys should be different
	same := true
	if len(key.Key) == len(key2.Key) {
		for i, b := range key.Key {
			if b != key2.Key[i] {
				same = false
				break
			}
		}
	} else {
		same = false
	}

	if same {
		t.Error("Different derivation paths should produce different keys")
	}
}

func TestSecureZeroMemory(t *testing.T) {
	// Create test data
	data := make([]byte, 100)
	for i := range data {
		data[i] = byte(i)
	}

	// Verify data is not zero
	allZero := true
	for _, b := range data {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("Test data should not be all zeros initially")
	}

	// Clear memory
	SecureZeroMemory(data)

	// Verify data is now zero
	for i, b := range data {
		if b != 0 {
			t.Errorf("Byte at index %d is not zero after SecureZeroMemory: %v", i, b)
		}
	}
}

func TestSecureRandom(t *testing.T) {
	sizes := []int{16, 32, 64, 128}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			data1, err := SecureRandom(size)
			if err != nil {
				t.Fatalf("Failed to generate random data: %v", err)
			}
			defer SecureZeroMemory(data1)

			if len(data1) != size {
				t.Errorf("Expected %d bytes, got %d", size, len(data1))
			}

			data2, err := SecureRandom(size)
			if err != nil {
				t.Fatalf("Failed to generate second random data: %v", err)
			}
			defer SecureZeroMemory(data2)

			// Data should be different (extremely high probability)
			same := true
			for i, b := range data1 {
				if b != data2[i] {
					same = false
					break
				}
			}
			if same {
				t.Error("Two random generations should produce different results")
			}
		})
	}
}

// Test invalid entropy sizes
func TestInvalidEntropySizes(t *testing.T) {
	invalidSizes := []int{64, 96, 300, 0, -1}

	for _, size := range invalidSizes {
		t.Run(fmt.Sprintf("invalid_%d", size), func(t *testing.T) {
			_, err := GenerateMnemonic(size)
			if err == nil {
				t.Errorf("Invalid entropy size accepted: %d", size)
			}
		})
	}
}

// Benchmark tests
func BenchmarkMnemonicGeneration128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mnemonic, err := GenerateMnemonic(128)
		if err != nil {
			b.Fatalf("Failed to generate mnemonic: %v", err)
		}
		_ = mnemonic
	}
}

func BenchmarkMnemonicGeneration256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mnemonic, err := GenerateMnemonic(256)
		if err != nil {
			b.Fatalf("Failed to generate mnemonic: %v", err)
		}
		_ = mnemonic
	}
}

func BenchmarkMnemonicToSeed(b *testing.B) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seed, err := MnemonicToSeed(mnemonic, "")
		if err != nil {
			b.Fatalf("Failed to generate seed: %v", err)
		}
		SecureZeroMemory(seed)
	}
}

func BenchmarkKeyDerivation(b *testing.B) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed, _ := MnemonicToSeed(mnemonic, "")
	defer SecureZeroMemory(seed)
	path := "m/44'/0'/0'/0/0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DeriveKey(seed, path)
		if err != nil {
			b.Fatalf("Failed to derive key: %v", err)
		}
	}
}

func BenchmarkSecureZeroMemory(b *testing.B) {
	data := make([]byte, 1024) // 1KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reinitialize data for each iteration
		for j := range data {
			data[j] = byte(j % 256)
		}
		SecureZeroMemory(data)
	}
}
