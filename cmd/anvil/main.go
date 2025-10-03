package main

import (
	"fmt"
	"os"
	"runtime"


	"anvil/internal/crypto"
	"anvil/pkg/wallet"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var (
	version = "dev" // Set via build flags
	buildDate = "unknown"
	gitCommit = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "anvil",
	Short:   "Anvil - Multi-cryptocurrency cold wallet generator",
	Long:    `A secure, offline multi-cryptocurrency cold wallet generator built in Go.`,
	Version: version,
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new wallet with multiple cryptocurrency accounts",
	Long: `Generate a new wallet using a BIP39 mnemonic phrase with accounts for
multiple cryptocurrencies. The generated wallet includes Bitcoin, Ethereum,
Dogecoin, and BNB accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateWallet()
	},
}

// recoverCmd represents the recover command
var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover a wallet from a mnemonic phrase",
	Long: `Recover a wallet from an existing BIP39 mnemonic phrase and generate
accounts for multiple cryptocurrencies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if mnemonic == "" {
			return fmt.Errorf("mnemonic phrase is required")
		}
		return recoverWallet(mnemonic)
	},
}

// deriveCmd represents the derive command  
var deriveCmd = &cobra.Command{
	Use:   "derive",
	Short: "Derive a specific account from a mnemonic",
	Long: `Derive a specific cryptocurrency account from a mnemonic phrase
using a custom derivation path.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if mnemonic == "" {
			return fmt.Errorf("mnemonic phrase is required")
		}
		if coinType == "" {
			return fmt.Errorf("coin type is required")
		}
		if path == "" {
			return fmt.Errorf("derivation path is required")
		}
		return deriveAccount(mnemonic, coinType, path)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(recoverCmd)
	rootCmd.AddCommand(deriveCmd)

	// Generate command flags
	generateCmd.Flags().IntVar(&words, "words", 12, "Number of words in mnemonic (12, 15, 18, 21, 24)")
	generateCmd.Flags().StringVar(&passphrase, "passphrase", "", "Optional passphrase for seed derivation")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for wallet data (default: stdout)")
	generateCmd.Flags().BoolVar(&includePrivate, "include-private", false, "Include private keys in output (DANGEROUS)")
	generateCmd.Flags().BoolVar(&includeMnemonic, "include-mnemonic", false, "Include mnemonic phrase in output (DANGEROUS)")
	generateCmd.Flags().StringVar(&format, "format", "json", "Output format: json, text, paper, qr")
	generateCmd.Flags().BoolVar(&paper, "paper", false, "Generate paper wallet format")
	generateCmd.Flags().BoolVar(&qrCodes, "qr", false, "Generate QR codes")

	// Recover command flags
	recoverCmd.Flags().StringVar(&mnemonic, "mnemonic", "", "BIP39 mnemonic phrase")
	recoverCmd.Flags().StringVar(&passphrase, "passphrase", "", "Optional passphrase for seed derivation")
	recoverCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for wallet data (default: stdout)")

	// Derive command flags
	deriveCmd.Flags().StringVar(&mnemonic, "mnemonic", "", "BIP39 mnemonic phrase")
	deriveCmd.Flags().StringVar(&coinType, "coin", "", "Coin type (BTC, ETH, DOGE, BNB)")
	deriveCmd.Flags().StringVar(&path, "path", "", "Derivation path (e.g., m/44'/0'/0'/0/0)")
	deriveCmd.Flags().StringVar(&passphrase, "passphrase", "", "Optional passphrase for seed derivation")
	deriveCmd.Flags().BoolVar(&includePrivate, "include-private", false, "Include private keys in output (DANGEROUS)")
	deriveCmd.Flags().StringVar(&format, "format", "json", "Output format: json, text")
	deriveCmd.Flags().BoolVar(&qrCodes, "qr", false, "Generate QR codes")
	deriveCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
}

func generateWallet() error {
	// Convert words to entropy bits
	entropyBits := wordsToEntropyBits(words)
	if entropyBits == 0 {
		return fmt.Errorf("invalid number of words: %d (must be 12, 15, 18, 21, or 24)", words)
	}

	// Generate mnemonic
	mnemonic, err := crypto.GenerateMnemonic(entropyBits)
	if err != nil {
		return fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return createWalletFromMnemonic(mnemonic, passphrase)
}

func recoverWallet(mnemonic string) error {
	return createWalletFromMnemonic(mnemonic, passphrase)
}

func createWalletFromMnemonic(mnemonic, passphrase string) error {
	// Convert mnemonic to seed
	seed, err := crypto.MnemonicToSeed(mnemonic, passphrase)
	if err != nil {
		return fmt.Errorf("failed to convert mnemonic to seed: %w", err)
	}
	defer crypto.ClearBytes(seed)

	// Create wallet with multiple coins
	wallet := &types.Wallet{
		Mnemonic:  mnemonic,
		Seed:      seed,
		Accounts:  []types.Account{},
		CoinTypes: make(map[string][]uint32),
		CreatedAt: time.Now(),
		Version:   version,
	}

	// Define coins to generate
	coins := []types.Coin{
		bitcoin.NewBitcoin(),
		ethereum.NewEthereum(),
		bitcoin.NewDogecoin(),
		ethereum.NewBinanceCoin(),
		tron.NewTron(),
	}

	// Generate accounts for each coin
	for _, coin := range coins {
		paths := getStandardPaths(coin)
		for _, path := range paths {
			account, err := coin.DeriveAccount(seed, path)
			if err != nil {
				return fmt.Errorf("failed to derive %s account: %w", coin.Symbol(), err)
			}
			wallet.Accounts = append(wallet.Accounts, account)
		}
		wallet.CoinTypes[coin.Symbol()] = []uint32{getCoinType(coin)}
	}

	// Output wallet data
	return outputWallet(wallet)
}

func deriveAccount(mnemonic, coinType, path string) error {
	// Convert mnemonic to seed
	seed, err := crypto.MnemonicToSeed(mnemonic, passphrase)
	if err != nil {
		return fmt.Errorf("failed to convert mnemonic to seed: %w", err)
	}
	defer crypto.ClearBytes(seed)

	// Get coin instance
	coin := getCoinInstance(coinType)
	if coin == nil {
		return fmt.Errorf("unsupported coin type: %s", coinType)
	}

	// Derive account
	account, err := coin.DeriveAccount(seed, path)
	if err != nil {
		return fmt.Errorf("failed to derive account: %w", err)
	}

	// Create minimal wallet for output
	wallet := &types.Wallet{
		Accounts: []types.Account{account},
		Version:  version,
	}

	return outputWallet(wallet)
}

func outputWallet(wallet *types.Wallet) error {
	// Determine output format
	outputFormat := types.OutputJSON
	if paper {
		outputFormat = types.OutputPaper
		includeMnemonic = true // Paper wallets should include mnemonic
	} else if qrCodes {
		outputFormat = types.OutputQR
	} else {
		switch format {
		case "json":
			outputFormat = types.OutputJSON
		case "text":
			outputFormat = types.OutputText
		case "paper":
			outputFormat = types.OutputPaper
			includeMnemonic = true
		case "qr":
			outputFormat = types.OutputQR
		default:
			return fmt.Errorf("unsupported format: %s", format)
		}
	}

	// Create output options
	options := types.OutputOptions{
		Format:          outputFormat,
		IncludePrivate:  includePrivate,
		IncludeMnemonic: includeMnemonic,
		FilePath:        outputFile,
	}

	// Validate options and show warnings
	if err := output.ValidateOptions(options); err != nil {
		return err
	}

	// Generate output
	generator := output.NewGenerator(options)
	return generator.GenerateWallet(wallet)
}
func wordsToEntropyBits(words int) int {
	switch words {
	case 12:
		return 128
	case 15:
		return 160
	case 18:
		return 192
	case 21:
		return 224
	case 24:
		return 256
	default:
		return 0
	}
}

func getStandardPaths(coin types.Coin) []string {
	switch c := coin.(type) {
	case *bitcoin.BitcoinCoin:
		return c.GetStandardDerivationPaths()
	case *ethereum.EthereumCoin:
		return c.GetStandardDerivationPaths()
	case *tron.TronCoin:
		return c.GetStandardDerivationPaths()
	default:
		// Default to BIP44 path
		return []string{"m/44'/0'/0'/0/0"}
	}
}

func getCoinType(coin types.Coin) uint32 {
	switch c := coin.(type) {
	case *bitcoin.BitcoinCoin:
		return c.GetCoinType()
	case *ethereum.EthereumCoin:
		return c.GetCoinType()
	case *tron.TronCoin:
		return c.GetCoinType()
	default:
		return 0
	}
}

func getCoinInstance(coinType string) types.Coin {
	switch coinType {
	case "BTC":
		return bitcoin.NewBitcoin()
	case "ETH":
		return ethereum.NewEthereum()
	case "DOGE":
		return bitcoin.NewDogecoin()
	case "BNB":
		return ethereum.NewBinanceCoin()
	case "TRX":
		return tron.NewTron()
	default:
		return nil
	}
}


func init() {
	// Initialize secure runtime on startup
	crypto.InitSecureRuntime()

	// Check for security warnings
	if warnings := crypto.VerifySecureEnvironment(); len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "⚠️  SECURITY WARNINGS:\n")
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "  • %s\n", warning)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Additional CLI flags for output formatting
var (
	version = "dev" // Set via build flags
	buildDate = "unknown"
	gitCommit = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Anvil Cold Wallet Generator\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Date: %s\n", buildDate)
		fmt.Printf("Git Commit: %s\n", gitCommit)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
