#!/bin/bash


# Function to build for Windows
build_windows() {
    echo "Building for Windows..."
    GOOS=windows GOARCH=amd64 go build -o "${PWD##*/}.exe"
    if [ $? -eq 0 ]; then
        echo "Build successful!"
    else
        echo "Build failed!"
    fi
}

# Function to build for Linux
build_linux() {
    echo "Building for Linux..."
    GOOS=linux GOARCH=amd64 go build -o "${PWD##*/}"
    if [ $? -eq 0 ]; then
        echo "Build successful!"
    else
        echo "Build failed!"
    fi
}

if [ "$1" == "win" ]; then
    build_windows
elif [ "$1" == "linux" ]; then
    build_linux
else
    echo "Invalid option. Usage: $0 [win|linux]"
    exit 1
fi
