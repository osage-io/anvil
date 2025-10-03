package bitcoin

import (
	"fmt"
	"time"

	"anvil/internal/crypto"
	"anvil/pkg/types"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// BitcoinCoin implements the types.Coin interface for Bitcoin
type BitcoinCoin struct {
	name      string
	symbol    string
	coinType  uint32
	netParams *chaincfg.Params
}

// NewBitcoin creates a new Bitcoin coin instance
func NewBitcoin() *BitcoinCoin {
	return &BitcoinCoin{
		name:      "Bitcoin",
		symbol:    "BTC",
		coinType:  0, // BIP44 coin type for Bitcoin
		netParams: &chaincfg.MainNetParams,
	}
}

// NewDogecoin creates a new Dogecoin instance (Bitcoin-like)
func NewDogecoin() *BitcoinCoin {
	// Create custom Dogecoin parameters based on Bitcoin mainnet
	dogecoinParams := chaincfg.MainNetParams
	dogecoinParams.Name = "dogecoin"
	dogecoinParams.Net = 0xc0c0c0c0
	dogecoinParams.PubKeyHashAddrID = 0x1E // Dogecoin addresses start with 'D'
	dogecoinParams.ScriptHashAddrID = 0x16 // P2SH addresses start with '9' or 'A'

	return &BitcoinCoin{
		name:      "Dogecoin",
		symbol:    "DOGE",
		coinType:  3, // BIP44 coin type for Dogecoin
		netParams: &dogecoinParams,
	}
}

// Name returns the full name of the cryptocurrency
func (b *BitcoinCoin) Name() string {
	return b.name
}

// Symbol returns the symbol/ticker of the cryptocurrency
func (b *BitcoinCoin) Symbol() string {
	return b.symbol
}

// DeriveAccount derives a new account for the given seed and derivation path
func (b *BitcoinCoin) DeriveAccount(seed []byte, path string) (types.Account, error) {
	// Derive the private key using BIP32
	key, err := crypto.DeriveKey(seed, path)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to derive key: %w", err)
	}

	// Get the private key bytes
	privateKeyBytes := key.Key

	// Create ECDSA private key from bytes
	privateKey, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	publicKeyBytes := publicKey.SerializeCompressed()

	// Generate address from public key
	address, err := b.publicKeyToAddress(publicKeyBytes)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to generate address: %w", err)
	}

	account := types.Account{
		Path:       path,
		PrivateKey: privateKey.Serialize(),
		PublicKey:  publicKeyBytes,
		Address:    address,
		Symbol:     b.symbol,
		CreatedAt:  time.Now(),
	}

	// Clear sensitive key data
	crypto.SecureZeroMemory(privateKeyBytes)

	return account, nil
}

// publicKeyToAddress converts a compressed public key to a Bitcoin address
func (b *BitcoinCoin) publicKeyToAddress(publicKeyBytes []byte) (string, error) {
	// Create address from compressed public key
	pubKeyHash := btcutil.Hash160(publicKeyBytes)
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, b.netParams)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %w", err)
	}

	return addr.EncodeAddress(), nil
}

// PrivateKeyToWIF converts a private key to Wallet Import Format
func (b *BitcoinCoin) PrivateKeyToWIF(privateKeyBytes []byte) (string, error) {
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	wif, err := btcutil.NewWIF(privateKey, b.netParams, true) // compressed
	if err != nil {
		return "", fmt.Errorf("failed to create WIF: %w", err)
	}
	return wif.String(), nil
}

// ValidateAddress checks if an address is valid for this cryptocurrency
func (b *BitcoinCoin) ValidateAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, b.netParams)
	return err == nil
}

// GetStandardDerivationPaths returns common derivation paths for this coin
func (b *BitcoinCoin) GetStandardDerivationPaths() []string {
	coinType := b.coinType
	return []string{
		fmt.Sprintf("m/44'/%d'/0'/0/0", coinType), // BIP44 (Legacy)
		fmt.Sprintf("m/49'/%d'/0'/0/0", coinType), // BIP49 (P2SH-P2WPKH)
		fmt.Sprintf("m/84'/%d'/0'/0/0", coinType), // BIP84 (Native SegWit)
	}
}

// GetCoinType returns the BIP44 coin type for this cryptocurrency
func (b *BitcoinCoin) GetCoinType() uint32 {
	return b.coinType
}
