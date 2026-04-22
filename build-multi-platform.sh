#!/bin/bash

# Build script for multi-platform release

echo "Building Necromancy for multiple platforms..."

# Create release directory
mkdir -p release

# Build for different platforms
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r goos goarch <<< "$platform"
    
    output_name="necromancy-${goos}-${goarch}"
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for $goos/$goarch..."
    
    env GOOS="$goos" GOARCH="$goarch" go build -ldflags="-s -w" -o "release/$output_name" .
    
    if [ $? -eq 0 ]; then
        echo "✅ Built: release/$output_name"
    else
        echo "❌ Failed to build for $goos/$goarch"
    fi
done

echo "Build complete! Check the release/ directory."

# Create checksums
cd release
shasum -a 256 * > checksums.txt
cd ..

echo "Checksums generated: release/checksums.txt"