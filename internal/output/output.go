package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"anvil/pkg/types"
	"github.com/skip2/go-qrcode"
)

// Generator handles different output formats for wallet data
type Generator struct {
	options types.OutputOptions
}

// NewGenerator creates a new output generator
func NewGenerator(options types.OutputOptions) *Generator {
	return &Generator{
		options: options,
	}
}

// GenerateWallet outputs wallet data in the specified format
func (g *Generator) GenerateWallet(wallet *types.Wallet) error {
	switch g.options.Format {
	case types.OutputJSON:
		return g.generateJSON(wallet)
	case types.OutputText:
		return g.generateText(wallet)
	case types.OutputPaper:
		return g.generatePaperWallet(wallet)
	case types.OutputQR:
		return g.generateQRCodes(wallet)
	default:
		return fmt.Errorf("unsupported output format")
	}
}

// generateJSON outputs wallet data as JSON
func (g *Generator) generateJSON(wallet *types.Wallet) error {
	var data []byte
	var err error

	if g.options.IncludePrivate || g.options.IncludeMnemonic {
		// Include sensitive data
		sensitive := struct {
			*types.Wallet
			IncludesPrivateKeys bool `json:"includes_private_keys"`
			IncludesMnemonic    bool `json:"includes_mnemonic"`
		}{
			Wallet:              wallet,
			IncludesPrivateKeys: g.options.IncludePrivate,
			IncludesMnemonic:    g.options.IncludeMnemonic,
		}

		if !g.options.IncludeMnemonic {
			sensitive.Wallet.Mnemonic = ""
			sensitive.Wallet.Seed = nil
		}

		if !g.options.IncludePrivate {
			// Clear private keys from accounts
			for i := range sensitive.Wallet.Accounts {
				sensitive.Wallet.Accounts[i].PrivateKey = nil
			}
		}

		data, err = json.MarshalIndent(sensitive, "", "  ")
	} else {
		// Safe JSON without sensitive data
		data, err = wallet.MarshalSafeJSON()
		if err == nil {
			// Pretty format the JSON
			var prettyJSON map[string]interface{}
			json.Unmarshal(data, &prettyJSON)
			data, err = json.MarshalIndent(prettyJSON, "", "  ")
		}
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return g.writeOutput(data, "json")
}

// generateText outputs wallet data as formatted text
func (g *Generator) generateText(wallet *types.Wallet) error {
	tmpl := `ANVIL CRYPTOCURRENCY WALLET
Generated: {{.CreatedAt.Format "2006-01-02 15:04:05"}}
Version: {{.Version}}

{{if .IncludeMnemonic}}MNEMONIC PHRASE (KEEP SECURE):
{{.Mnemonic}}

{{end}}ACCOUNTS:
{{range .Accounts}}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{{.Symbol}} - {{.Path}}
Address: {{.Address}}
{{if $.IncludePrivate}}Private Key: {{printf "%x" .PrivateKey}}{{end}}
Public Key:  {{printf "%x" .PublicKey}}
Created: {{.CreatedAt.Format "2006-01-02 15:04:05"}}

{{end}}
SUPPORTED CRYPTOCURRENCIES:
{{range $symbol, $types := .CoinTypes}}• {{$symbol}} (BIP44 Type: {{index $types 0}})
{{end}}

⚠️  SECURITY WARNING ⚠️
• Keep this information secure and offline
• Never share your private keys or mnemonic phrase
• Store multiple backups in secure locations
• Verify addresses before sending funds
`

	t, err := template.New("wallet").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		*types.Wallet
		IncludeMnemonic bool
		IncludePrivate  bool
	}{
		Wallet:          wallet,
		IncludeMnemonic: g.options.IncludeMnemonic,
		IncludePrivate:  g.options.IncludePrivate,
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return g.writeOutput([]byte(buf.String()), "txt")
}

// generatePaperWallet outputs wallet data in paper wallet format
func (g *Generator) generatePaperWallet(wallet *types.Wallet) error {
	tmpl := `
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃                           ANVIL PAPER WALLET                            ┃
┃                          🔨 KEEP THIS SECURE 🔨                         ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

Generated: {{.CreatedAt.Format "January 2, 2006 at 15:04:05"}}

{{if .IncludeMnemonic}}
┌─ RECOVERY PHRASE (BIP39) ─────────────────────────────────────────────┐
│  {{.Mnemonic}}                                                         
└───────────────────────────────────────────────────────────────────────┘
{{end}}

{{range .Accounts}}
┌─ {{.Symbol}} WALLET ──────────────────────────────────────────────────────┐
│ Path: {{.Path}}
│ Address: {{.Address}}
{{if $.IncludePrivate}}│ Private: {{printf "%x" .PrivateKey}}{{end}}
└───────────────────────────────────────────────────────────────────────┘

{{end}}
┌─ INSTRUCTIONS ────────────────────────────────────────────────────────┐
│ 1. Keep this paper wallet in a secure, dry location                   │
│ 2. Make multiple copies and store in different locations              │
│ 3. Never share your private keys or recovery phrase                   │
│ 4. Use the recovery phrase to restore wallets in compatible software  │
│ 5. Always verify addresses before sending funds                       │
└───────────────────────────────────────────────────────────────────────┘

⚠️  This paper contains sensitive cryptographic keys. Treat it like cash! ⚠️
`

	t, err := template.New("paper").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		*types.Wallet
		IncludeMnemonic bool
		IncludePrivate  bool
	}{
		Wallet:          wallet,
		IncludeMnemonic: g.options.IncludeMnemonic,
		IncludePrivate:  g.options.IncludePrivate,
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return g.writeOutput([]byte(buf.String()), "txt")
}

// generateQRCodes generates QR codes for wallet data
func (g *Generator) generateQRCodes(wallet *types.Wallet) error {
	baseDir := g.options.FilePath
	if baseDir == "" {
		baseDir = "anvil-qr-codes"
	}

	// Remove extension if present
	if ext := filepath.Ext(baseDir); ext != "" {
		baseDir = strings.TrimSuffix(baseDir, ext)
	}

	// Create directory
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create QR codes directory: %w", err)
	}

	// Generate mnemonic QR code if requested
	if g.options.IncludeMnemonic && wallet.Mnemonic != "" {
		qrFile := filepath.Join(baseDir, "mnemonic.png")
		if err := qrcode.WriteFile(wallet.Mnemonic, qrcode.Medium, 256, qrFile); err != nil {
			return fmt.Errorf("failed to generate mnemonic QR code: %w", err)
		}
		fmt.Printf("Generated mnemonic QR code: %s\n", qrFile)
	}

	// Generate QR codes for each account
	for _, account := range wallet.Accounts {
		prefix := fmt.Sprintf("%s-%s", strings.ToLower(account.Symbol),
			strings.ReplaceAll(account.Path, "/", "-")[2:]) // Remove m/ prefix

		// Address QR code
		addrFile := filepath.Join(baseDir, prefix+"-address.png")
		if err := qrcode.WriteFile(account.Address, qrcode.Medium, 256, addrFile); err != nil {
			return fmt.Errorf("failed to generate address QR code: %w", err)
		}

		// Private key QR code if requested
		if g.options.IncludePrivate && len(account.PrivateKey) > 0 {
			privFile := filepath.Join(baseDir, prefix+"-private.png")
			privHex := fmt.Sprintf("%x", account.PrivateKey)
			if err := qrcode.WriteFile(privHex, qrcode.Medium, 256, privFile); err != nil {
				return fmt.Errorf("failed to generate private key QR code: %w", err)
			}
		}

		fmt.Printf("Generated QR codes for %s %s\n", account.Symbol, account.Address)
	}

	// Generate info file
	infoFile := filepath.Join(baseDir, "README.txt")
	info := fmt.Sprintf(`ANVIL QR CODES
Generated: %s

This directory contains QR codes for your cryptocurrency wallet.

FILES:
`, time.Now().Format("2006-01-02 15:04:05"))

	if g.options.IncludeMnemonic {
		info += "• mnemonic.png - Your BIP39 recovery phrase (KEEP SECURE!)\n"
	}

	for _, account := range wallet.Accounts {
		prefix := fmt.Sprintf("%s-%s", strings.ToLower(account.Symbol),
			strings.ReplaceAll(account.Path, "/", "-")[2:])
		info += fmt.Sprintf("• %s-address.png - %s address: %s\n", prefix, account.Symbol, account.Address)
		if g.options.IncludePrivate {
			info += fmt.Sprintf("• %s-private.png - %s private key (KEEP SECURE!)\n", prefix, account.Symbol)
		}
	}

	info += `
SECURITY WARNINGS:
• QR codes containing private keys or mnemonic phrases are extremely sensitive
• Store these files on secure, offline devices only
• Delete private key QR codes after use
• Never share or upload these files online
• Make backup copies on separate secure media
`

	if err := os.WriteFile(infoFile, []byte(info), 0644); err != nil {
		return fmt.Errorf("failed to write info file: %w", err)
	}

	fmt.Printf("QR codes saved to directory: %s\n", baseDir)
	return nil
}

// writeOutput writes data to file or stdout
func (g *Generator) writeOutput(data []byte, defaultExt string) error {
	if g.options.FilePath == "" {
		// Write to stdout
		fmt.Print(string(data))
		return nil
	}

	// Ensure file has proper extension
	filePath := g.options.FilePath
	if filepath.Ext(filePath) == "" {
		filePath += "." + defaultExt
	}

	return os.WriteFile(filePath, data, 0644)
}

// ValidateOptions checks if the output options are valid
func ValidateOptions(options types.OutputOptions) error {
	// Warn about sensitive data
	if options.IncludePrivate {
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: Output will include private keys! Ensure secure handling.\n")
	}

	if options.IncludeMnemonic {
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: Output will include mnemonic phrase! Ensure secure handling.\n")
	}

	// Check file path for QR codes
	if options.Format == types.OutputQR && options.FilePath != "" {
		// Ensure directory exists or can be created
		dir := filepath.Dir(options.FilePath)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("cannot create output directory: %w", err)
			}
		}
	}

	return nil
}
