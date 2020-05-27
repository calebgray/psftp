#!/bin/bash

# Install Prerequisites
add-apt-repository ppa:longsleep/golang-backports
apt update
apt upgrade -y
apt install -y build-essential golang-go curl git-all

# Windows Shim
OUTPUT_SUFFIX=""
if [ "$GOOS" == 'windows' ]; then
  OUTPUT_SUFFIX='.exe'
  apt install -y mingw-w64
  [ "$GOARCH" == '386' ] && CCARCH=i686 || CCARCH=x86_64
  export CC=${CCARCH}-w64-mingw32-gcc
fi

# Set the Output Binary Name
OUTPUT_NAME="${PROJECT_NAME}-${GOARCH}-${PROJECT_VERSION}"

# Setup Go Build Environment
PROJECT_ROOT="/go/src/github.com/${GITHUB_REPOSITORY}"
PROJECT_PARENT_DIR="/go/src/github.com/${GITHUB_ACTOR}"
mkdir -p "$PROJECT_PARENT_DIR"
ln -s "$GITHUB_WORKSPACE" "$PROJECT_ROOT"
cd "$PROJECT_ROOT" || exit

# Fetch Dependencies
go get -v ./...

# Run the Build
go build . -o "${OUTPUT_NAME}${OUTPUT_SUFFIX}"

# Create the Archive
zip -r9 "${OUTPUT_NAME}.zip" "${OUTPUT_NAME}${OUTPUT_SUFFIX}"

# Upload the Release
curl "${UPLOAD_URL}?name=${OUTPUT_NAME}.zip" -X POST --data-binary "@${OUTPUT_NAME}.zip" -H 'Content-Type: application/zip' -H "Authorization: Bearer ${GITHUB_TOKEN}"
