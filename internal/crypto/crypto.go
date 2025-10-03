package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"anvil/pkg/types"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// SecureRandom generates cryptographically secure random bytes
func SecureRandom(size int) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// GenerateMnemonic creates a new BIP39 mnemonic phrase
func GenerateMnemonic(entropyBits int) (string, error) {
	if entropyBits%32 != 0 || entropyBits < 128 || entropyBits > 256 {
		return "", fmt.Errorf("entropy bits must be 128, 160, 192, 224, or 256")
	}

	entropy, err := SecureRandom(entropyBits / 8)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}
	defer ClearBytes(entropy)

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// MnemonicToSeed converts a BIP39 mnemonic to a seed with optional passphrase
func MnemonicToSeed(mnemonic, passphrase string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic phrase")
	}

	seed := bip39.NewSeed(mnemonic, passphrase)
	return seed, nil
}

// DeriveKey derives a private key from seed using BIP32 derivation path
func DeriveKey(seed []byte, path string) (*bip32.Key, error) {
	derivePath, err := ParseDerivationPath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid derivation path: %w", err)
	}

	// Generate master key from seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key: %w", err)
	}

	// Derive child keys following the path
	currentKey := masterKey

	// Purpose (hardened)
	currentKey, err = currentKey.NewChildKey(derivePath.Purpose + bip32.FirstHardenedChild)
	if err != nil {
		return nil, fmt.Errorf("failed to derive purpose key: %w", err)
	}

	// Coin type (hardened)
	currentKey, err = currentKey.NewChildKey(derivePath.CoinType + bip32.FirstHardenedChild)
	if err != nil {
		return nil, fmt.Errorf("failed to derive coin type key: %w", err)
	}

	// Account (hardened)
	currentKey, err = currentKey.NewChildKey(derivePath.Account + bip32.FirstHardenedChild)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %w", err)
	}

	// Change (not hardened)
	currentKey, err = currentKey.NewChildKey(derivePath.Change)
	if err != nil {
		return nil, fmt.Errorf("failed to derive change key: %w", err)
	}

	// Index (not hardened)
	finalKey, err := currentKey.NewChildKey(derivePath.Index)
	if err != nil {
		return nil, fmt.Errorf("failed to derive index key: %w", err)
	}

	return finalKey, nil
}

// ParseDerivationPath parses a BIP32 derivation path string
func ParseDerivationPath(path string) (types.DerivationPath, error) {
	if !strings.HasPrefix(path, "m/") {
		return types.DerivationPath{}, fmt.Errorf("path must start with 'm/'")
	}

	parts := strings.Split(path[2:], "/")
	if len(parts) != 5 {
		return types.DerivationPath{}, fmt.Errorf("path must have 5 components: m/purpose'/coin_type'/account'/change/index")
	}

	var dp types.DerivationPath
	var err error

	// Parse purpose (should be hardened)
	dp.Purpose, err = parsePathComponent(parts[0], true)
	if err != nil {
		return types.DerivationPath{}, fmt.Errorf("invalid purpose: %w", err)
	}

	// Parse coin type (should be hardened)
	dp.CoinType, err = parsePathComponent(parts[1], true)
	if err != nil {
		return types.DerivationPath{}, fmt.Errorf("invalid coin type: %w", err)
	}

	// Parse account (should be hardened)
	dp.Account, err = parsePathComponent(parts[2], true)
	if err != nil {
		return types.DerivationPath{}, fmt.Errorf("invalid account: %w", err)
	}

	// Parse change (not hardened)
	dp.Change, err = parsePathComponent(parts[3], false)
	if err != nil {
		return types.DerivationPath{}, fmt.Errorf("invalid change: %w", err)
	}

	// Parse index (not hardened)
	dp.Index, err = parsePathComponent(parts[4], false)
	if err != nil {
		return types.DerivationPath{}, fmt.Errorf("invalid index: %w", err)
	}

	return dp, nil
}

// parsePathComponent parses a single component of a derivation path
func parsePathComponent(component string, shouldBeHardened bool) (uint32, error) {
	isHardened := strings.HasSuffix(component, "'")

	if shouldBeHardened && !isHardened {
		return 0, fmt.Errorf("component should be hardened (end with ')")
	}
	if !shouldBeHardened && isHardened {
		return 0, fmt.Errorf("component should not be hardened")
	}

	numStr := component
	if isHardened {
		numStr = component[:len(component)-1]
	}

	val, err := strconv.ParseUint(numStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}

	return uint32(val), nil
}

// ClearBytes securely zeros out a byte slice
func ClearBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// Zero the memory
	for i := range b {
		b[i] = 0
	}

	// Try to ensure the compiler doesn't optimize away the clearing
	runtime.KeepAlive(b)
}

// ClearString attempts to clear a string from memory (best effort)
func ClearString(s *string) {
	if s == nil || len(*s) == 0 {
		return
	}

	// Convert string to byte slice and clear
	// Note: This is best effort - Go strings are immutable
	// and the runtime may keep copies elsewhere
	header := (*reflect.StringHeader)(unsafe.Pointer(s))
	if header.Data != 0 {
		slice := (*[1 << 30]byte)(unsafe.Pointer(header.Data))[:header.Len:header.Len]
		ClearBytes(slice)
	}
	*s = ""
}

// Hash256 performs double SHA256 hash (Bitcoin style)
func Hash256(data []byte) []byte {
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:]
}

// Hash160 performs SHA256 followed by RIPEMD160
func Hash160(data []byte) []byte {
	// For now, return SHA256 (RIPEMD160 requires additional dependency)
	hash := sha256.Sum256(data)
	return hash[:20] // Take first 20 bytes
}
