#!/usr/bin/env bash

# Anvil Installation Script
# Automatically downloads and installs the latest release for your platform

set -euo pipefail

# Configuration
REPO="osage-io/anvil"
INSTALL_DIR="$HOME/bin"
BINARY_NAME="anvil"
GITHUB_API="https://api.github.com"
GITHUB_RELEASES="https://github.com"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${GREEN}[INFO]${NC} $*"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

debug() {
    if [[ "${DEBUG:-}" == "1" ]]; then
        echo -e "${BLUE}[DEBUG]${NC} $*" >&2
    fi
}

# Detect platform and architecture
detect_platform() {
    local os arch
    
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    
    case "$os" in
        linux)
            case "$arch" in
                x86_64|amd64)
                    echo "linux-amd64"
                    ;;
                aarch64|arm64)
                    echo "linux-arm64"
                    ;;
                *)
                    error "Unsupported architecture: $arch"
                    exit 1
                    ;;
            esac
            ;;
        darwin)
            case "$arch" in
                x86_64|amd64)
                    echo "darwin-amd64"
                    ;;
                arm64)
                    echo "darwin-arm64"
                    ;;
                *)
                    error "Unsupported architecture: $arch"
                    exit 1
                    ;;
            esac
            ;;
        mingw*|msys*|cygwin*|windows*)
            case "$arch" in
                x86_64|amd64)
                    echo "windows-amd64"
                    ;;
                aarch64|arm64)
                    echo "windows-arm64"
                    ;;
                *)
                    error "Unsupported architecture: $arch"
                    exit 1
                    ;;
            esac
            ;;
        *)
            error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Get latest release info from GitHub API
get_latest_release() {
    local url="$GITHUB_API/repos/$REPO/releases/latest"
    
    debug "Fetching latest release info from: $url"
    
    if command_exists curl; then
        curl -s "$url"
    elif command_exists wget; then
        wget -q -O- "$url"
    else
        error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Download file with progress
download_file() {
    local url="$1"
    local output="$2"
    
    log "Downloading from: $url"
    
    if command_exists curl; then
        curl -L --progress-bar "$url" -o "$output"
    elif command_exists wget; then
        wget --progress=bar:force -O "$output" "$url"
    else
        error "Neither curl nor wget is available"
        exit 1
    fi
}

# Verify checksum if available
verify_checksum() {
    local file="$1"
    local expected_checksum="$2"
    
    if [[ -z "$expected_checksum" ]]; then
        warn "No checksum provided, skipping verification"
        return 0
    fi
    
    if command_exists sha256sum; then
        local actual_checksum
        actual_checksum=$(sha256sum "$file" | cut -d' ' -f1)
        
        if [[ "$actual_checksum" == "$expected_checksum" ]]; then
            log "Checksum verification passed"
            return 0
        else
            error "Checksum verification failed!"
            error "Expected: $expected_checksum"
            error "Actual:   $actual_checksum"
            return 1
        fi
    else
        warn "sha256sum not available, skipping checksum verification"
        return 0
    fi
}

# Extract archive
extract_archive() {
    local archive="$1"
    local dest_dir="$2"
    
    debug "Extracting $archive to $dest_dir"
    
    case "$archive" in
        *.tar.gz|*.tgz)
            tar -xzf "$archive" -C "$dest_dir"
            ;;
        *.zip)
            if command_exists unzip; then
                unzip -q "$archive" -d "$dest_dir"
            else
                error "unzip command not found. Please install unzip to extract .zip files."
                exit 1
            fi
            ;;
        *)
            error "Unsupported archive format: $archive"
            exit 1
            ;;
    esac
}

# Main installation function
install_anvil() {
    local platform version download_url archive_name temp_dir
    local checksum_url checksums expected_checksum=""
    
    log "🔨 Anvil Installation Script"
    log "=============================="
    
    # Detect platform
    platform=$(detect_platform)
    log "Detected platform: $platform"
    
    # Get latest release info
    log "Fetching latest release information..."
    local release_json
    release_json=$(get_latest_release)
    
    # Parse version
    version=$(echo "$release_json" | grep -o '"tag_name": *"[^"]*"' | sed 's/"tag_name": *"\(.*\)"/\1/')
    if [[ -z "$version" ]]; then
        error "Failed to determine latest version"
        exit 1
    fi
    log "Latest version: $version"
    
    # Construct download URL and archive name
    archive_name="anvil-${version}-${platform}"
    
    # Determine archive extension
    if [[ "$platform" == windows-* ]]; then
        archive_name="${archive_name}.zip"
    else
        archive_name="${archive_name}.tar.gz"
    fi
    
    download_url="$GITHUB_RELEASES/$REPO/releases/download/$version/$archive_name"
    
    debug "Archive name: $archive_name"
    debug "Download URL: $download_url"
    
    # Create temporary directory
    temp_dir=$(mktemp -d)
    debug "Using temporary directory: $temp_dir"
    
    # Cleanup function
    cleanup() {
        debug "Cleaning up temporary files..."
        rm -rf "$temp_dir"
    }
    trap cleanup EXIT
    
    # Download the release archive
    local archive_path="$temp_dir/$archive_name"
    download_file "$download_url" "$archive_path"
    
    # Try to download and verify checksums
    checksum_url="$GITHUB_RELEASES/$REPO/releases/download/$version/checksums.txt"
    checksums="$temp_dir/checksums.txt"
    
    if download_file "$checksum_url" "$checksums" 2>/dev/null; then
        expected_checksum=$(grep "$archive_name" "$checksums" 2>/dev/null | cut -d' ' -f1 || echo "")
        debug "Expected checksum: $expected_checksum"
    else
        warn "Could not download checksums file, skipping verification"
    fi
    
    # Verify checksum if available
    if [[ -n "$expected_checksum" ]]; then
        verify_checksum "$archive_path" "$expected_checksum"
    fi
    
    # Create installation directory
    if [[ ! -d "$INSTALL_DIR" ]]; then
        log "Creating installation directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi
    
    # Extract archive
    log "Extracting archive..."
    extract_archive "$archive_path" "$temp_dir"
    
    # Find the binary in extracted files
    local binary_path
    if [[ "$platform" == windows-* ]]; then
        binary_path=$(find "$temp_dir" -name "${BINARY_NAME}.exe" | head -1)
    else
        binary_path=$(find "$temp_dir" -name "$BINARY_NAME" -type f | head -1)
    fi
    
    if [[ -z "$binary_path" ]]; then
        error "Could not find $BINARY_NAME binary in extracted files"
        ls -la "$temp_dir"
        exit 1
    fi
    
    debug "Found binary at: $binary_path"
    
    # Install binary
    local target_path="$INSTALL_DIR/$BINARY_NAME"
    if [[ "$platform" == windows-* ]]; then
        target_path="${target_path}.exe"
    fi
    
    log "Installing to: $target_path"
    cp "$binary_path" "$target_path"
    chmod +x "$target_path"
    
    # Verify installation
    log "Verifying installation..."
    if [[ -x "$target_path" ]]; then
        local installed_version
        installed_version=$("$target_path" version 2>/dev/null | grep "Version:" | awk '{print $2}' || echo "unknown")
        log "✅ Successfully installed Anvil $installed_version"
        
        # Check if install dir is in PATH
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            warn "⚠️  $INSTALL_DIR is not in your PATH"
            warn "   Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            warn "   export PATH=\"\$PATH:$INSTALL_DIR\""
            warn ""
            warn "   Or run: echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc"
        else
            log "✅ $INSTALL_DIR is already in your PATH"
        fi
        
        log ""
        log "🎉 Installation complete!"
        log "   Run 'anvil --help' to get started"
        log "   Run 'anvil version' to verify the installation"
    else
        error "Installation verification failed"
        exit 1
    fi
}

# Help function
show_help() {
    cat << HELP_END
Anvil Installation Script

USAGE:
    $0 [OPTIONS]

OPTIONS:
    -h, --help      Show this help message
    -d, --debug     Enable debug output
    -v, --version   Show script version

EXAMPLES:
    # Install latest version
    $0

    # Install with debug output
    DEBUG=1 $0

ENVIRONMENT VARIABLES:
    INSTALL_DIR     Installation directory (default: $HOME/bin)
    DEBUG          Enable debug output (set to 1)

For more information, visit: https://github.com/$REPO
HELP_END
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -d|--debug)
            export DEBUG=1
            shift
            ;;
        -v|--version)
            echo "Anvil installation script v1.0.0"
            exit 0
            ;;
        *)
            error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Run main installation
install_anvil
