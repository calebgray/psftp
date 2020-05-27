#!/bin/bash

# Install Prerequisites
apt update
apt upgrade -y
apt install -y build-essential curl git-all

# Install Go
cd /
wget https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz
tar xf go1.14.3.linux-amd64.tar.gz
export PATH=$PATH:/go/bin

# Windows Shim
export OUTPUT_SUFFIX=""
if [ "$GOOS" == 'windows' ]; then
	OUTPUT_SUFFIX='.exe'
	apt install -y mingw-w64
	[ "$GOARCH" == '386' ] && CCARCH=i686 || CCARCH=x86_64
	export CC=${CCARCH}-w64-mingw32-gcc
fi

# Set the Output Binary Name
export OUTPUT_NAME="${PROJECT_NAME}-${GOARCH}-${PROJECT_VERSION}"

# Setup Go Build Environment
PROJECT_ROOT="/go/src/github.com/${GITHUB_REPOSITORY}"
PROJECT_PARENT_DIR="/go/src/github.com/${GITHUB_ACTOR}"
mkdir -p "$PROJECT_PARENT_DIR"
ln -s "$GITHUB_WORKSPACE" "$PROJECT_ROOT"
cd "$PROJECT_ROOT" || exit

# Fetch Dependencies
go get -v ./...

# Run the Build
go build -o "${OUTPUT_NAME}${OUTPUT_SUFFIX}" .

# Create the Archive
zip -r9 "${OUTPUT_NAME}.zip" "${OUTPUT_NAME}${OUTPUT_SUFFIX}"
