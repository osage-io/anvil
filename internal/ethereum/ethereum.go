package ethereum

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"anvil/internal/crypto"
	"anvil/pkg/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// EthereumCoin implements the types.Coin interface for Ethereum-based coins
type EthereumCoin struct {
	name     string
	symbol   string
	coinType uint32
	chainID  uint64
}

// NewEthereum creates a new Ethereum coin instance
func NewEthereum() *EthereumCoin {
	return &EthereumCoin{
		name:     "Ethereum",
		symbol:   "ETH",
		coinType: 60, // BIP44 coin type for Ethereum
		chainID:  1,  // Ethereum Mainnet
	}
}

// NewBinanceCoin creates a new Binance Smart Chain coin instance
func NewBinanceCoin() *EthereumCoin {
	return &EthereumCoin{
		name:     "BNB Smart Chain",
		symbol:   "BNB",
		coinType: 60, // Uses same coin type as Ethereum
		chainID:  56, // BSC Mainnet
	}
}

// Name returns the full name of the cryptocurrency
func (e *EthereumCoin) Name() string {
	return e.name
}

// Symbol returns the symbol/ticker of the cryptocurrency
func (e *EthereumCoin) Symbol() string {
	return e.symbol
}

// DeriveAccount derives a new account for the given seed and derivation path
func (e *EthereumCoin) DeriveAccount(seed []byte, path string) (types.Account, error) {
	// Derive the private key using BIP32
	key, err := crypto.DeriveKey(seed, path)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to derive key: %w", err)
	}

	// Get the private key bytes
	privateKeyBytes := key.Key

	// Create ECDSA private key from bytes
	privateKey, err := ethcrypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to create ECDSA key: %w", err)
	}

	// Get uncompressed public key bytes (65 bytes)
	publicKeyBytes := ethcrypto.FromECDSAPub(&privateKey.PublicKey)

	// Generate Ethereum address from public key
	address := e.publicKeyToAddress(publicKeyBytes)

	account := types.Account{
		Path:       path,
		PrivateKey: privateKeyBytes,
		PublicKey:  publicKeyBytes,
		Address:    address,
		Symbol:     e.symbol,
		CreatedAt:  time.Now(),
	}

	// Clear sensitive key data
	crypto.SecureZeroMemory(privateKeyBytes)

	return account, nil
}

// publicKeyToAddress converts an uncompressed public key to an Ethereum address
func (e *EthereumCoin) publicKeyToAddress(publicKeyBytes []byte) string {
	// Remove the 0x04 prefix if present (uncompressed key indicator)
	if len(publicKeyBytes) == 65 && publicKeyBytes[0] == 0x04 {
		publicKeyBytes = publicKeyBytes[1:]
	}

	// Hash the public key with Keccak256
	hash := ethcrypto.Keccak256Hash(publicKeyBytes)

	// Take the last 20 bytes as the address
	address := common.BytesToAddress(hash[12:])

	// Return checksummed address (EIP-55)
	return e.toChecksumAddress(address.Hex())
}

// ValidateAddress checks if an address is a valid Ethereum address
func (e *EthereumCoin) ValidateAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}

	// Check if it's a valid hex address
	if !common.IsHexAddress(address) {
		return false
	}

	// If it has mixed case, verify EIP-55 checksum
	if e.hasMixedCase(address) {
		return e.isValidChecksum(address)
	}

	return true
}

// hasMixedCase checks if an address has mixed case letters
func (e *EthereumCoin) hasMixedCase(address string) bool {
	address = strings.TrimPrefix(address, "0x")
	hasUpper := false
	hasLower := false

	for _, char := range address {
		if char >= 'A' && char <= 'F' {
			hasUpper = true
		} else if char >= 'a' && char <= 'f' {
			hasLower = true
		}
	}

	return hasUpper && hasLower
}

// isValidChecksum verifies EIP-55 checksum
func (e *EthereumCoin) isValidChecksum(address string) bool {
	return address == e.toChecksumAddress(address)
}

// GetStandardDerivationPaths returns common derivation paths for Ethereum
func (e *EthereumCoin) GetStandardDerivationPaths() []string {
	coinType := e.coinType
	return []string{
		fmt.Sprintf("m/44'/%d'/0'/0/0", coinType), // BIP44 standard path
		fmt.Sprintf("m/44'/%d'/0'/0/1", coinType), // Second address
		fmt.Sprintf("m/44'/%d'/1'/0/0", coinType), // Change addresses
	}
}

// GetCoinType returns the BIP44 coin type for this cryptocurrency
func (e *EthereumCoin) GetCoinType() uint32 {
	return e.coinType
}

// GetChainID returns the chain ID for this network
func (e *EthereumCoin) GetChainID() uint64 {
	return e.chainID
}

// toChecksumAddress converts an address to EIP-55 checksum format
func (e *EthereumCoin) toChecksumAddress(address string) string {
	address = strings.ToLower(strings.TrimPrefix(address, "0x"))
	hash := ethcrypto.Keccak256Hash([]byte(address))

	result := "0x"
	hashHex := hex.EncodeToString(hash[:])

	for i, char := range address {
		if char >= '0' && char <= '9' {
			result += string(char)
		} else {
			hashChar := hashHex[i]
			if hashChar >= '8' {
				result += strings.ToUpper(string(char))
			} else {
				result += string(char)
			}
		}
	}

	return result
}
