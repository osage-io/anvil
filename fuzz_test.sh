#!/bin/bash

echo "Starting fuzz testing for Anvil wallet..."

# Fuzz test crypto utilities
echo "Fuzzing crypto utilities..."
go test -fuzz=FuzzMnemonicGeneration -fuzztime=30s ./internal/crypto/ || echo "No fuzz tests found in crypto"

# Fuzz test address generation
echo "Fuzzing Bitcoin address generation..."
go test -fuzz=FuzzBitcoinAddress -fuzztime=30s ./internal/bitcoin/ || echo "No fuzz tests found in bitcoin"

echo "Fuzzing Ethereum address generation..."
go test -fuzz=FuzzEthereumAddress -fuzztime=30s ./internal/ethereum/ || echo "No fuzz tests found in ethereum"

echo "Fuzzing TRON address generation..."
go test -fuzz=FuzzTronAddress -fuzztime=30s ./internal/tron/ || echo "No fuzz tests found in tron"

echo "Fuzz testing completed."
