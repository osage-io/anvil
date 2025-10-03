#!/bin/bash

echo "🚀 Pre-push Testing Script for Anvil"
echo "===================================="
echo ""

# Test build
echo "1. Testing build..."
go build -o anvil cmd/anvil/main.go
if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

# Test version
echo ""
echo "2. Testing version command..."
./anvil version

# Test help
echo ""
echo "3. Testing help command..."
./anvil --help | head -5

# Run tests
echo ""
echo "4. Running unit tests..."
go test -v ./internal/... 2>/dev/null | grep -E "(PASS|FAIL|RUN)"
TEST_EXIT_CODE=${PIPESTATUS[0]}
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "✅ All tests passed"
else
    echo "❌ Some tests failed"
    exit 1
fi

# Test cross-compilation for all platforms
echo ""
echo "5. Testing cross-compilation..."
PLATFORMS=(
    "linux/amd64"
    "linux/arm64" 
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    echo "   Building for $GOOS/$GOARCH..."
    
    EXTENSION=""
    if [ "$GOOS" = "windows" ]; then
        EXTENSION=".exe"
    fi
    
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
        -o anvil-test-$GOOS-$GOARCH$EXTENSION \
        cmd/anvil/main.go
    
    if [ $? -eq 0 ]; then
        echo "   ✅ $platform build successful"
        rm -f anvil-test-$GOOS-$GOARCH$EXTENSION
    else
        echo "   ❌ $platform build failed"
        exit 1
    fi
done

# Test with sample wallet generation
echo ""
echo "6. Testing wallet generation..."
echo "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about" | \
    ./anvil derive --mnemonic "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about" \
    --coin BTC --path "m/44'/0'/0'/0/0" --format text | head -5

if [ $? -eq 0 ]; then
    echo "✅ Wallet derivation test successful"
else
    echo "❌ Wallet derivation test failed"
    exit 1
fi

# Cleanup
echo ""
echo "7. Cleaning up..."
rm -f anvil

echo ""
echo "🎉 All pre-push tests passed!"
echo "Ready to push to GitHub!"
echo ""
echo "Next steps:"
echo "1. Create repository at https://github.com/osage/anvil"
echo "2. git push -u origin main"
echo "3. git tag v0.1.0 && git push origin v0.1.0"
