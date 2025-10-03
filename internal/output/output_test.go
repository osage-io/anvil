package output

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"anvil/pkg/types"
)

func TestNewGenerator(t *testing.T) {
	options := types.OutputOptions{
		Format:          types.OutputJSON,
		IncludePrivate:  false,
		IncludeMnemonic: false,
		FilePath:        "",
	}

	generator := NewGenerator(options)
	if generator == nil {
		t.Fatal("NewGenerator returned nil")
	}
}

func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name    string
		options types.OutputOptions
		wantErr bool
	}{
		{
			name: "valid json options",
			options: types.OutputOptions{
				Format:          types.OutputJSON,
				IncludePrivate:  false,
				IncludeMnemonic: false,
				FilePath:        "",
			},
			wantErr: false,
		},
		{
			name: "valid text options",
			options: types.OutputOptions{
				Format:          types.OutputText,
				IncludePrivate:  false,
				IncludeMnemonic: false,
				FilePath:        "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createTestWallet() *types.Wallet {
	return &types.Wallet{
		Version:   "test",
		CreatedAt: time.Now(),
		Accounts: []types.Account{
			{
				Path:       "m/44'/0'/0'/0/0",
				PrivateKey: []byte("test-private-key"),
				PublicKey:  []byte("0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798"),
				Address:    "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				Symbol:     "BTC",
				CreatedAt:  time.Now(),
			},
		},
		CoinTypes: map[string][]uint32{"BTC": {0}},
	}
}

func TestGenerateJSON(t *testing.T) {
	wallet := createTestWallet()

	generator := &Generator{
		options: types.OutputOptions{
			Format:          types.OutputJSON,
			IncludePrivate:  false,
			IncludeMnemonic: false,
			FilePath:        "",
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := generator.generateJSON(wallet)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("generateJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "BTC") {
		t.Error("JSON output should contain BTC")
	}
	if !strings.Contains(output, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa") {
		t.Error("JSON output should contain address")
	}
}

func TestGenerateText(t *testing.T) {
	wallet := createTestWallet()

	generator := &Generator{
		options: types.OutputOptions{
			Format:          types.OutputText,
			IncludePrivate:  false,
			IncludeMnemonic: false,
			FilePath:        "",
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := generator.generateText(wallet)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("generateText failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ANVIL CRYPTOCURRENCY WALLET") {
		t.Error("Text output should contain wallet header")
	}
	if !strings.Contains(output, "BTC") {
		t.Error("Text output should contain BTC")
	}
}

func TestWriteOutput(t *testing.T) {
	generator := &Generator{
		options: types.OutputOptions{
			FilePath: "",
		},
	}

	data := []byte("test output")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := generator.writeOutput(data, "txt")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("writeOutput failed: %v", err)
	}

	if buf.String() != "test output" {
		t.Errorf("Expected 'test output', got %s", buf.String())
	}
}

func TestGenerateWallet(t *testing.T) {
	wallet := createTestWallet()

	generator := NewGenerator(types.OutputOptions{
		Format:          types.OutputJSON,
		IncludePrivate:  false,
		IncludeMnemonic: false,
		FilePath:        "",
	})

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := generator.GenerateWallet(wallet)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("GenerateWallet failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "BTC") {
		t.Error("Output should contain BTC")
	}
}

func TestGenerateWalletWithText(t *testing.T) {
	wallet := createTestWallet()

	generator := NewGenerator(types.OutputOptions{
		Format:          types.OutputText,
		IncludePrivate:  false,
		IncludeMnemonic: false,
		FilePath:        "",
	})

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := generator.GenerateWallet(wallet)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("GenerateWallet failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ANVIL CRYPTOCURRENCY WALLET") {
		t.Error("Text output should contain wallet header")
	}
}

func TestGenerateWalletUnsupportedFormat(t *testing.T) {
	wallet := createTestWallet()

	generator := &Generator{
		options: types.OutputOptions{
			Format: 999, // Invalid format
		},
	}

	err := generator.GenerateWallet(wallet)
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}
