# Anvil Installation Scripts

This directory contains installation and utility scripts for Anvil.

## install.sh

An automated installation script that downloads and installs the latest Anvil release for your platform.

### Features

- **Cross-platform**: Supports Linux, macOS, and Windows
- **Architecture detection**: Automatically detects amd64/arm64
- **Secure downloads**: Verifies checksums when available
- **Progress indicators**: Shows download progress
- **PATH integration**: Checks if installation directory is in PATH
- **Clean installation**: Extracts to temporary directory and cleans up

### Usage

```bash
# Install latest version
./scripts/install.sh

# Install with debug output
DEBUG=1 ./scripts/install.sh

# Show help
./scripts/install.sh --help
```

### One-liner Installation

You can install Anvil directly from GitHub with:

```bash
curl -fsSL https://raw.githubusercontent.com/osage-io/anvil/main/scripts/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/osage-io/anvil/main/scripts/install.sh | bash
```

### Environment Variables

- `INSTALL_DIR`: Installation directory (default: `$HOME/bin`)
- `DEBUG`: Enable debug output (set to `1`)

### Examples

```bash
# Install to a custom directory
INSTALL_DIR="$HOME/.local/bin" ./scripts/install.sh

# Install with debug output
DEBUG=1 ./scripts/install.sh

# Install to system-wide location (requires sudo)
sudo INSTALL_DIR="/usr/local/bin" ./scripts/install.sh
```

### Supported Platforms

- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64, arm64

### Dependencies

The script automatically detects and uses available tools:
- Download: `curl` or `wget`
- Extraction: `tar` (for .tar.gz), `unzip` (for .zip)
- Verification: `sha256sum` (optional)

### Security

- Downloads are verified with SHA256 checksums when available
- Uses HTTPS for all downloads
- Temporary files are cleaned up automatically
- Binary permissions are set correctly (`chmod +x`)

### Troubleshooting

1. **Download fails**: Ensure you have `curl` or `wget` installed
2. **Extraction fails**: Ensure you have `tar` and/or `unzip` installed
3. **Permission denied**: Check that `INSTALL_DIR` is writable
4. **Command not found**: Add `INSTALL_DIR` to your PATH

```bash
# Add to PATH permanently
echo 'export PATH="$PATH:$HOME/bin"' >> ~/.bashrc
source ~/.bashrc
```
