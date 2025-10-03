# CLI Reference

Complete reference for all Anvil CLI commands, flags, and options.

## Global Options

```
  -h, --help      help for anvil
  -v, --version   version for anvil
```

## Commands

### `anvil generate`

Generate a new wallet with multiple cryptocurrency accounts.

#### Synopsis

```bash
anvil generate [flags]
```

#### Description

Creates a new hierarchical deterministic (HD) wallet using a BIP39 mnemonic phrase. The wallet includes accounts for Bitcoin, Ethereum, Dogecoin, and BNB Smart Chain by default.

#### Examples

```bash
# Generate with default 12-word mnemonic
anvil generate

# Generate with 24-word mnemonic (more secure)
anvil generate --words 24

# Generate with passphrase protection
anvil generate --words 24 --passphrase "my-secure-passphrase"

# Save output to file
anvil generate --words 24 --output wallet.json

# Generate and show all derivation paths
anvil generate --words 12 --verbose
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--words` | int | 12 | Number of words in mnemonic (12, 15, 18, 21, 24) |
| `--passphrase` | string | "" | Optional passphrase for seed derivation (BIP39) |
| `--output, -o` | string | "" | Output file for wallet data (default: stdout) |

#### Output Format

The command outputs a JSON structure containing:

```json
{
  "accounts": [
    {
      "path": "m/44'/0'/0'/0/0",
      "public_key": "base64-encoded-public-key",
      "address": "cryptocurrency-address",
      "symbol": "BTC",
      "created_at": "2023-10-03T14:00:00Z"
    }
  ],
  "coin_types": {
    "BTC": [0],
    "ETH": [60],
    "DOGE": [3],
    "BNB": [60]
  },
  "created_at": "2023-10-03T14:00:00Z",
  "version": "0.1.0"
}
```

**Note**: Private keys and mnemonic are excluded from output for security. Use `--include-private` to include them (use with extreme caution).

---

### `anvil recover`

Recover a wallet from an existing mnemonic phrase.

#### Synopsis

```bash
anvil recover [flags]
```

#### Description

Recovers a hierarchical deterministic wallet from an existing BIP39 mnemonic phrase and generates accounts for all supported cryptocurrencies.

#### Examples

```bash
# Recover from 12-word mnemonic
anvil recover --mnemonic "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12"

# Recover with passphrase
anvil recover --mnemonic "your mnemonic phrase" --passphrase "your-passphrase"

# Recover and save to file
anvil recover --mnemonic "your mnemonic phrase" --output recovered-wallet.json
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--mnemonic` | string | **required** | BIP39 mnemonic phrase (12-24 words) |
| `--passphrase` | string | "" | Optional passphrase used during original generation |
| `--output, -o` | string | "" | Output file for wallet data (default: stdout) |

#### Mnemonic Validation

The command validates the mnemonic phrase according to BIP39 standards:
- Word count must be 12, 15, 18, 21, or 24
- All words must be from the BIP39 English wordlist
- Checksum must be valid

---

### `anvil derive`

Derive a specific cryptocurrency account from a mnemonic phrase.

#### Synopsis

```bash
anvil derive [flags]
```

#### Description

Derives a single cryptocurrency account using a specific derivation path. Useful for generating individual addresses or exploring different derivation paths.

#### Examples

```bash
# Derive Bitcoin SegWit address
anvil derive --coin BTC --path "m/84'/0'/0'/0/0" --mnemonic "your mnemonic"

# Derive Ethereum address
anvil derive --coin ETH --path "m/44'/60'/0'/0/0" --mnemonic "your mnemonic"

# Derive Dogecoin address with passphrase
anvil derive --coin DOGE --path "m/44'/3'/0'/0/5" --mnemonic "your mnemonic" --passphrase "passphrase"
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--mnemonic` | string | **required** | BIP39 mnemonic phrase |
| `--coin` | string | **required** | Coin type (BTC, ETH, DOGE, BNB, TRX) |
| `--path` | string | **required** | BIP44 derivation path (e.g., m/44'/0'/0'/0/0) |
| `--passphrase` | string | "" | Optional passphrase for seed derivation |

#### Supported Coins

| Coin | Symbol | BIP44 Type | Standard Paths |
|------|--------|------------|----------------|
| Bitcoin | BTC | 0 | m/44'/0'/0'/0/x, m/49'/0'/0'/0/x, m/84'/0'/0'/0/x |
| Ethereum | ETH | 60 | m/44'/60'/0'/0/x |
| Dogecoin | DOGE | 3 | m/44'/3'/0'/0/x, m/49'/3'/0'/0/x, m/84'/3'/0'/0/x |
| BNB Smart Chain | BNB | 60 | m/44'/60'/0'/0/x |
| TRON | TRX | 195 | m/44'/195'/0'/0/x |

#### Derivation Path Format

Derivation paths follow the BIP44 standard: `m/purpose'/coin_type'/account'/change/address_index`

- `purpose'`: Usually 44 (BIP44), 49 (BIP49 P2SH-P2WPKH), or 84 (BIP84 native SegWit)
- `coin_type'`: BIP44 registered coin type (hardened)
- `account'`: Account index (hardened, usually 0)
- `change`: 0 for receiving addresses, 1 for change addresses
- `address_index`: Sequential address index within the account

**Note**: The `'` indicates hardened derivation.

---

### `anvil version`

Display version information.

#### Synopsis

```bash
anvil version
```

#### Description

Shows the current version of Anvil along with build information.

---

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ANVIL_CONFIG_DIR` | Configuration directory | `~/.config/anvil` |
| `ANVIL_NO_COLOR` | Disable colored output | `false` |

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Cryptographic error |
| 4 | File I/O error |

## Security Considerations

### Safe Practices

1. **Use on Air-Gapped Systems**: Run Anvil on computers without network connectivity
2. **Verify Checksums**: Always verify the integrity of downloaded binaries
3. **Secure Storage**: Store mnemonic phrases in secure, offline locations
4. **Test First**: Use small amounts when testing generated addresses
5. **Multiple Backups**: Create multiple backups of your mnemonic phrase

### Unsafe Practices

❌ **Never** run Anvil on compromised or untrusted systems
❌ **Never** share your mnemonic phrase or private keys
❌ **Never** store sensitive data in plain text files
❌ **Never** use weak or predictable passphrases
❌ **Never** generate wallets on systems connected to the internet (if holding significant value)

## Troubleshooting

### Common Issues

**Error: "invalid mnemonic phrase"**
- Verify all words are from the BIP39 English wordlist
- Check for typos or extra spaces
- Ensure word count is 12, 15, 18, 21, or 24

**Error: "invalid derivation path"**
- Verify path format: `m/purpose'/coin_type'/account'/change/address_index`
- Ensure hardened components (account, coin_type, purpose) end with `'`
- Check that coin type matches the selected cryptocurrency

**Error: "unsupported coin type"**
- Use supported coin symbols: BTC, ETH, DOGE, BNB, TRX
- Check for typos in coin symbol

### Getting Help

- Check the [FAQ](../concepts/faq.md)
- Review [Security Best Practices](../security/)
- Report bugs on GitHub: [github.com/yourorg/anvil/issues](https://github.com/yourorg/anvil/issues)
