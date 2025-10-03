# Anvil CLI Documentation

Welcome to the Anvil CLI documentation. This guide covers all aspects of using Anvil, the multi-cryptocurrency cold wallet generator.

## Table of Contents

- [Quick Start](#quick-start)
- [CLI Reference](cli/)
- [Core Concepts](concepts/)
- [Security Model](security/)

## Quick Start

### Installation

```bash
# Install via Go (recommended)
go install github.com/yourorg/anvil/cmd/anvil@latest

# Or build from source
git clone https://github.com/yourorg/anvil.git
cd anvil
go build -o anvil ./cmd/anvil
```

### Generate Your First Wallet

```bash
# Generate a new 12-word wallet
anvil generate

# Generate a 24-word wallet (more secure)
anvil generate --words 24

# Save to file
anvil generate --words 24 --output my-wallet.json
```

### Recover from Existing Mnemonic

```bash
anvil recover --mnemonic "word1 word2 word3 ... word12"
```

### Derive Specific Addresses

```bash
# Bitcoin address
anvil derive --coin BTC --path "m/84'/0'/0'/0/0" --mnemonic "your mnemonic"

# Ethereum address  
anvil derive --coin ETH --path "m/44'/60'/0'/0/0" --mnemonic "your mnemonic"
```

## Supported Cryptocurrencies

| Cryptocurrency | Symbol | Description |
|---------------|--------|-------------|
| Bitcoin | BTC | Original cryptocurrency with multiple address formats |
| Ethereum | ETH | Smart contract platform with EIP-55 checksums |
| Dogecoin | DOGE | Bitcoin-based with custom address prefixes |
| BNB Smart Chain | BNB | Ethereum-compatible blockchain |
| TRON | TRX | High-throughput blockchain platform |

## Security Features

- ✅ **Offline Operation**: No network connectivity required
- ✅ **Memory Safety**: Sensitive data cleared after use
- ✅ **Industry Standards**: BIP39, BIP32, BIP44 compliant
- ✅ **Safe Defaults**: Private keys excluded from output by default
- ✅ **Cross-Platform**: Works on Linux, macOS, Windows

## Next Steps

- Read the [CLI Reference](cli/) for detailed command documentation
- Learn about [Core Concepts](concepts/) like HD wallets and derivation paths
- Review the [Security Model](security/) to understand Anvil's security approach
