package tron

import (
	"crypto/sha256"
	"fmt"
	"time"

	"anvil/internal/crypto"
	"anvil/pkg/types"
	"github.com/btcsuite/btcd/btcutil/base58"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// TronCoin implements the types.Coin interface for TRON
type TronCoin struct {
	name     string
	symbol   string
	coinType uint32
}

// NewTron creates a new TRON coin instance
func NewTron() *TronCoin {
	return &TronCoin{
		name:     "TRON",
		symbol:   "TRX",
		coinType: 195, // BIP44 coin type for TRON
	}
}

// Name returns the full name of the cryptocurrency
func (t *TronCoin) Name() string {
	return t.name
}

// Symbol returns the symbol/ticker of the cryptocurrency
func (t *TronCoin) Symbol() string {
	return t.symbol
}

// DeriveAccount derives a new account for the given seed and derivation path
func (t *TronCoin) DeriveAccount(seed []byte, path string) (types.Account, error) {
	// Derive the private key using BIP32
	key, err := crypto.DeriveKey(seed, path)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to derive key: %w", err)
	}
	
	// Get the private key bytes
	privateKeyBytes := key.Key
	
	// Create ECDSA private key from bytes (same as Ethereum)
	privateKey, err := ethcrypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to create ECDSA key: %w", err)
	}
	
	// Get uncompressed public key bytes (65 bytes)
	publicKeyBytes := ethcrypto.FromECDSAPub(&privateKey.PublicKey)
	
	// Generate TRON address from public key
	address, err := t.publicKeyToAddress(publicKeyBytes)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to generate address: %w", err)
	}
	
	account := types.Account{
		Path:       path,
		PrivateKey: privateKeyBytes,
		PublicKey:  publicKeyBytes,
		Address:    address,
		Symbol:     t.symbol,
		CreatedAt:  time.Now(),
	}
	
	// Clear sensitive key data
	crypto.SecureZeroMemory(privateKeyBytes)
	
	return account, nil
}

// publicKeyToAddress converts an uncompressed public key to a TRON address
func (t *TronCoin) publicKeyToAddress(publicKeyBytes []byte) (string, error) {
	// Remove the 0x04 prefix if present (uncompressed key indicator)
	if len(publicKeyBytes) == 65 && publicKeyBytes[0] == 0x04 {
		publicKeyBytes = publicKeyBytes[1:]
	}
	
	// Hash the public key with Keccak256 (same as Ethereum)
	hash := ethcrypto.Keccak256(publicKeyBytes)
	
	// Take the last 20 bytes as the address
	addressBytes := hash[12:]
	
	// Add TRON prefix (0x41) to make it 21 bytes
	tronAddressBytes := make([]byte, 21)
	tronAddressBytes[0] = 0x41 // TRON mainnet prefix
	copy(tronAddressBytes[1:], addressBytes)
	
	// TRON uses double SHA256 for checksum, not Bitcoin's CheckEncode
	checksum := doubleSHA256(tronAddressBytes)
	
	// Append first 4 bytes of checksum
	addressWithChecksum := make([]byte, 25)
	copy(addressWithChecksum[:21], tronAddressBytes)
	copy(addressWithChecksum[21:], checksum[:4])
	
	// Encode with Base58 (no version byte)
	address := base58.Encode(addressWithChecksum)
	
	return address, nil
}

// ValidateAddress checks if an address is a valid TRON address
func (t *TronCoin) ValidateAddress(address string) bool {
	// Decode the Base58 address
	decoded := base58.Decode(address)
	if len(decoded) != 25 {
		return false
	}
	
	// Split address and checksum
	addressBytes := decoded[:21]
	providedChecksum := decoded[21:]
	
	// Check TRON mainnet prefix
	if addressBytes[0] != 0x41 {
		return false
	}
	
	// Verify checksum
	expectedChecksum := doubleSHA256(addressBytes)
	for i := 0; i < 4; i++ {
		if providedChecksum[i] != expectedChecksum[i] {
			return false
		}
	}
	
	return true
}

// GetStandardDerivationPaths returns common derivation paths for TRON
func (t *TronCoin) GetStandardDerivationPaths() []string {
	coinType := t.coinType
	return []string{
		fmt.Sprintf("m/44'/%d'/0'/0/0", coinType),  // BIP44 standard path
		fmt.Sprintf("m/44'/%d'/0'/0/1", coinType),  // Second address
		fmt.Sprintf("m/44'/%d'/1'/0/0", coinType),  // Change addresses
	}
}

// GetCoinType returns the BIP44 coin type for TRON
func (t *TronCoin) GetCoinType() uint32 {
	return t.coinType
}

// AddressToHex converts a TRON Base58 address to hex format
func (t *TronCoin) AddressToHex(address string) (string, error) {
	// Decode the Base58 address
	decoded := base58.Decode(address)
	if len(decoded) != 25 {
		return "", fmt.Errorf("invalid address length")
	}
	
	// Extract the 21-byte address (skip checksum)
	addressBytes := decoded[:21]
	if addressBytes[0] != 0x41 {
		return "", fmt.Errorf("invalid TRON address prefix")
	}
	
	// Convert to hex (remove 0x41 prefix and return as 0x...)
	hexAddress := fmt.Sprintf("0x%x", addressBytes[1:])
	return hexAddress, nil
}

// HexToAddress converts a hex address to TRON Base58 format
func (t *TronCoin) HexToAddress(hexAddress string) (string, error) {
	// Remove 0x prefix if present
	if len(hexAddress) >= 2 && hexAddress[:2] == "0x" {
		hexAddress = hexAddress[2:]
	}
	
	// Check hex address length (should be 40 characters = 20 bytes)
	if len(hexAddress) != 40 {
		return "", fmt.Errorf("invalid hex address length")
	}
	
	// Parse hex string to bytes
	addressBytes := make([]byte, 20)
	for i := 0; i < 20; i++ {
		_, err := fmt.Sscanf(hexAddress[i*2:i*2+2], "%02x", &addressBytes[i])
		if err != nil {
			return "", fmt.Errorf("invalid hex address format: %w", err)
		}
	}
	
	// Add TRON prefix
	tronAddressBytes := make([]byte, 21)
	tronAddressBytes[0] = 0x41
	copy(tronAddressBytes[1:], addressBytes)
	
	// Add checksum
	checksum := doubleSHA256(tronAddressBytes)
	addressWithChecksum := make([]byte, 25)
	copy(addressWithChecksum[:21], tronAddressBytes)
	copy(addressWithChecksum[21:], checksum[:4])
	
	// Encode with Base58
	address := base58.Encode(addressWithChecksum)
	return address, nil
}

// doubleSHA256 performs double SHA256 hash (TRON checksum method)
func doubleSHA256(data []byte) []byte {
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:]
}
