# Solana Support

Anvil now supports generating secure, offline Solana (SOL) wallets using industry-standard BIP39 mnemonic phrases.

## Solana Key Derivation

Solana uses ED25519 cryptography with the following derivation path pattern:

- **Standard path**: `m/44'/501'/0'/0'` (account 0)
- **Multiple accounts**: `m/44'/501'/1'/0'`, `m/44'/501'/2'/0'`, etc.
- **Coin type**: 501 (BIP44 registered)

All path components are hardened (use apostrophes) in Solana, which is different from Bitcoin-style coins where the last two components (change and address index) are typically non-hardened.

## Examples

### Generate a Multi-Coin Wallet (Including Solana)

```bash
anvil generate --words 12
```

This creates wallets for Bitcoin, Ethereum, Dogecoin, BNB, TRON, and **Solana**.

### Generate Only Solana Address

```bash
anvil derive --coin SOL --mnemonic "your twelve word phrase here" --path "m/44'/501'/0'/0'"
```

### Generate Multiple Solana Accounts

```bash
# First account
anvil derive --coin SOL --mnemonic "your phrase" --path "m/44'/501'/0'/0'"

# Second account  
anvil derive --coin SOL --mnemonic "your phrase" --path "m/44'/501'/1'/0'"

# Third account
anvil derive --coin SOL --mnemonic "your phrase" --path "m/44'/501'/2'/0'"
```

## Solana Address Format

Solana addresses are base58-encoded public keys (32 bytes), for example:
- `E7sEpbacn6HMfhkpR9Rqj8e6fCuunLJzNpRWb67sw1wE`
- `A1bxvMv5VDN3EAumTpSiZzSwnrPVgvFbw3aT932weUUN`

## Testing on Devnet

To test your generated addresses:

1. **Switch to Solana devnet**:
   ```bash
   solana config set --url https://api.devnet.solana.com
   ```

2. **Check balance**:
   ```bash
   solana balance <YOUR_ADDRESS>
   ```

3. **Request airdrop** (devnet only):
   ```bash
   solana airdrop 1 <YOUR_ADDRESS>
   ```

## Security Notes

- **Offline Generation**: Anvil generates all keys offline without any network connections
- **Mnemonic Security**: Store your mnemonic phrase securely - it can recover all your wallets
- **Private Keys**: Never share private keys or import them into online tools
- **Test First**: Always test with small amounts on devnet/testnet before using on mainnet

## Compatible Wallets

The generated Solana addresses and private keys are compatible with:
- Phantom Wallet
- Solflare Wallet
- Ledger Hardware Wallets
- solana-keygen CLI tool

## Network Support

Anvil generates addresses that work on:
- **Mainnet-beta** (production)
- **Testnet** (testing)
- **Devnet** (development)

The same address works across all networks - only the RPC endpoint changes.
