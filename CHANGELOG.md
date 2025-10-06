# Changelog

## [v0.2.0] - 2025-10-06

### Added
- **Solana (SOL) Support**: Added full support for generating Solana wallets
  - ED25519 key generation using proper HD derivation paths (m/44'/501'/X'/0')
  - Base58-encoded addresses compatible with all major Solana wallets
  - Support for mainnet, testnet, and devnet
  - Comprehensive unit tests and documentation
- Support for deriving individual Solana accounts using `--coin SOL`
- Documentation in `docs/solana.md` with usage examples

### Changed
- Updated CLI help text to include SOL in supported coin types
- Added Solana to default multi-coin wallet generation
- README updated with Solana support information
