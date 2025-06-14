#!/bin/bash

# Define your application name and main package path
APP_NAME="konfigo"
PACKAGE_PATH="./cmd/konfigo" # Or "cmd/myapp/main.go" or "./cmd/myapp"

# Define your target platforms (OS/ARCH)
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
    "freebsd/amd64"
    "freebsd/arm64"
    "openbsd/amd64"
    "openbsd/arm64"
    "netbsd/amd64"
    "netbsd/arm64"
)

# Create a directory for the builds
BUILD_DIR="builds"
mkdir -p "$BUILD_DIR"

for platform in "${PLATFORMS[@]}"
do
    # Split the platform string into OS and ARCH
    IFS='/' read -r GOOS GOARCH <<< "$platform"

    # Define the output name
    OUTPUT_NAME="$BUILD_DIR/${APP_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT_NAME" "$PACKAGE_PATH"
    if [ $? -ne 0 ]; then
        echo "An error occurred while building for $GOOS/$GOARCH."
        # Optionally exit on error: exit 1
    fi
done

echo "All builds completed."