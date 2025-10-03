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

```bash
go install github.com/yourorg/anvil@latest
```

## Quick Start

```bash
# Generate a new wallet
anvil generate

# Generate with custom options
anvil generate --words 24 --paper --out wallet.json

# Recover from existing mnemonic
anvil recover --mnemonic "your twelve word seed phrase here..."

# Derive specific addresses
anvil derive --coin BTC --path "m/84'/0'/0'/0/0"
```

## Supported Cryptocurrencies

- Bitcoin (BTC)
- Ethereum (ETH) 
- Binance Coin (BNB)
- Dogecoin (DOGE)
- TRON (TRX)
- More coming soon...

## License

MIT License - see LICENSE file for details.
