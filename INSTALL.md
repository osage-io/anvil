# Anvil Cold Wallet Generator - Installation Guide

## Quick Installation

### Download Pre-built Binaries (Recommended)

1. Go to the [Releases page](https://github.com/osage/anvil/releases)
2. Download the appropriate binary for your platform:
   - **Linux**: `anvil-linux-amd64.tar.gz` or `anvil-linux-arm64.tar.gz`
   - **macOS**: `anvil-darwin-amd64.tar.gz` (Intel) or `anvil-darwin-arm64.tar.gz` (Apple Silicon)
   - **Windows**: `anvil-windows-amd64.zip` or `anvil-windows-arm64.zip`

3. Extract the archive:
   ```bash
   # Linux/macOS
   tar -xzf anvil-*.tar.gz
   
   # Windows (use your preferred extraction tool)
   unzip anvil-windows-*.zip
   ```

4. Make executable (Linux/macOS):
   ```bash
   chmod +x anvil-*
   ```

5. Optionally, move to PATH:
   ```bash
   # Linux/macOS
   sudo mv anvil-* /usr/local/bin/anvil
   
   # Or add to your PATH
   export PATH=$PATH:/path/to/anvil
   ```

### Verify Installation

```bash
# Check version
./anvil version

# Get help
./anvil --help
```

### Security Verification

Always verify the integrity of downloaded binaries:

```bash
# Download checksums.txt from the release
# Verify SHA256
sha256sum -c checksums.txt
```

## Build from Source

### Prerequisites

- Go 1.21 or later
- Git

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/osage/anvil.git
   cd anvil
   ```

2. Build:
   ```bash
   go build -o anvil cmd/anvil/main.go
   ```

3. Run tests:
   ```bash
   go test -v ./internal/...
   ```

## First Use

### Generate a New Wallet
```bash
# Generate with default settings
./anvil generate

# Generate with specific entropy and format
./anvil generate --entropy-size 256 --format text
```

### Derive Specific Account
```bash
./anvil derive \
  --mnemonic "your twelve word mnemonic phrase here" \
  --coin BTC \
  --path "m/44'/0'/0'/0/0" \
  --format text
```

## Security Recommendations

1. **Air-Gapped Machine**: Run Anvil on a computer that has never been connected to the internet
2. **Verify Downloads**: Always check SHA256 checksums
3. **Clean Environment**: Use a fresh OS installation
4. **Secure Storage**: Store mnemonic phrases in multiple secure locations
5. **Test First**: Always test with small amounts before using for large sums

## Supported Platforms

- ✅ Linux (x86_64, ARM64)
- ✅ macOS (Intel, Apple Silicon) 
- ✅ Windows (x86_64, ARM64)

## Need Help?

- 📖 Read the [full documentation](./WARP.md)
- 🧪 Check [testing guide](./TESTING.md)
- 🐛 Report issues on [GitHub](https://github.com/osage/anvil/issues)
