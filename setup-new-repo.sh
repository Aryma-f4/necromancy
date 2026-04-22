#!/bin/bash

# Setup script for new Necromancy repository

echo "🚀 Setting up Necromancy repository..."

# Check if git is available
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed. Please install Git first."
    exit 1
fi

# Initialize git repository
echo "📁 Initializing Git repository..."
git init

# Create .gitignore if not exists
if [ ! -f .gitignore ]; then
    echo "📝 Creating .gitignore..."
    cat > .gitignore << 'EOF'
# Binaries
necromancy
necromancy-*
*.exe

# Logs
*.log
logs/

# Go build cache
*.o
*.a
*.so

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Dependency directories
vendor/

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# Release directory
release/
dist/
EOF
fi

# Add all files
echo "📂 Adding files to Git..."
git add .

# Create initial commit
echo "💾 Creating initial commit..."
git commit -m "Initial commit: Necromancy Advanced Shell Manager

- Multi-platform post-exploitation tool
- Complete feature set from penelope.py
- Enhanced with ASCII banner and utilities
- CI/CD ready for automated releases

Features:
- Reverse shell management
- Session multiplexing
- Post-exploitation modules
- File transfer capabilities
- Port forwarding
- Multi-platform support (Linux, macOS, Windows)"

# Add remote repository (user needs to replace with their repo)
echo "🔗 Repository setup complete!"
echo ""
echo "To push to GitHub:"
echo "1. Create a new repository at: https://github.com/new"
echo "2. Add your remote: git remote add origin https://github.com/YOUR_USERNAME/necromancy.git"
echo "3. Push: git push -u origin main"
echo ""
echo "Or push to the target repo:"
echo "git remote add origin https://github.com/Aryma-f4/necromancy.git"
echo "git push -u origin main"
echo ""
echo "✅ Setup complete! Ready for CI/CD and multi-platform releases."