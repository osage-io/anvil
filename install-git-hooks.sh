#!/usr/bin/env bash
set -e

echo "🔧 Installing Git hooks..."

# Get the repository root and hooks directory
REPO_ROOT="$(git rev-parse --show-toplevel)"
HOOK_SRC_DIR="$REPO_ROOT/.githooks"
HOOK_DEST_DIR="$(git rev-parse --git-path hooks)"

# Check if the .githooks directory exists
if [ ! -d "$HOOK_SRC_DIR" ]; then
    echo "❌ .githooks directory not found at $HOOK_SRC_DIR"
    exit 1
fi

# Install each hook
for hook in pre-commit pre-push; do
    src_file="$HOOK_SRC_DIR/$hook"
    dest_file="$HOOK_DEST_DIR/$hook"
    
    if [ -f "$src_file" ]; then
        echo "  • Installing $hook hook..."
        
        # Make source file executable
        chmod +x "$src_file"
        
        # Create symlink (or copy if symlink fails)
        if ln -sf "$src_file" "$dest_file" 2>/dev/null; then
            echo "    ✅ Symlinked $hook"
        else
            # Fallback to copying if symlinking fails (e.g., on Windows)
            cp "$src_file" "$dest_file"
            chmod +x "$dest_file"
            echo "    ✅ Copied $hook (symlink not supported)"
        fi
    else
        echo "    ⚠️  $hook hook not found, skipping"
    fi
done

echo ""
echo "✅ Git hooks installation complete!"
echo ""
echo "📋 What's installed:"
echo "   • pre-commit: Runs formatting, linting, and fast tests"
echo "   • pre-push: Runs comprehensive test suite"
echo ""
echo "💡 Tips:"
echo "   • To bypass hooks in emergencies: git commit --no-verify"
echo "   • To bypass push hooks: git push --no-verify"
echo "   • To uninstall: rm .git/hooks/pre-commit .git/hooks/pre-push"
echo ""
