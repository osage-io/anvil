# Anvil 🔨

A secure, offline multi-cryptocurrency cold wallet generator built in Go.

## Features

- **Multi-coin support**: Bitcoin, Ethereum, Dogecoin, BNB, TRON, and more
- **Offline security**: No network connectivity required or used
- **BIP39 mnemonics**: Industry-standard seed phrase generation
- **Multiple output formats**: JSON, text, QR codes, paper wallets
- **Hierarchical Deterministic (HD) wallets**: BIP44/49/84 derivation paths
- **Cross-platform**: Linux, macOS, Windows

## Security First

Anvil is designed for maximum security:
- Memory is cleared after use
- No network connections
- Cryptographically secure random number generation
- Open source and auditable

## Installation

### Quick Install (Recommended)

```bash
# Install latest release
curl -fsSL https://raw.githubusercontent.com/osage-io/anvil/main/scripts/install.sh | bash
```

This will automatically download and install the latest release to `~/bin`.

### Manual Installation

```bash
go install github.com/yourorg/anvil@latest
```

## Quick Start

```bash
# Generate a new wallet
anvil generate

# Generate 24-word wallet and save to file
anvil generate --words 24 --output ~/crypto/test_wallet.json

# Generate with custom options
anvil generate --words 24 --paper --out wallet.json

# Recover from existing mnemonic
anvil recover --mnemonic "your twelve word seed phrase here..."

# Derive specific addresses
anvil derive --coin SOL --path "m/44'/501'/0'/0'"
anvil derive --coin BTC --path "m/84'/0'/0'/0/0"
```

## Supported Cryptocurrencies

- Bitcoin (BTC)
- Ethereum (ETH) 
- Binance Coin (BNB)
- Dogecoin (DOGE)
- TRON (TRX)
- Solana (SOL)
- More coming soon...

## License

MIT License - see LICENSE file for details.
