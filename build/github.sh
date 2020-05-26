#!/bin/sh

# Set the Output Binary Name
[ -z "${PROJECT_VERSION}" ] && exit 1
OUTPUT_NAME="${PROJECT_NAME}-${GOARCH}-${PROJECT_VERSION}"

# Setup Go Build Environment
PROJECT_ROOT="/go/src/github.com/${GITHUB_REPOSITORY}"
PROJECT_PARENT_DIR=$(dirname "$PROJECT_ROOT")
mkdir -p "$PROJECT_PARENT_DIR"
ln -s "$GITHUB_WORKSPACE" "$PROJECT_ROOT"
cd "$PROJECT_ROOT" || exit

# Fetch Dependencies
go get -v ./...

# Run the Build
OUTPUT_SUFFIX=""
[ "$GOOS" == 'windows' ] && OUTPUT_SUFFIX='.exe'
go build . -o "${OUTPUT_NAME}${OUTPUT_SUFFIX}"

# Create the Archive
zip -r9 "${OUTPUT_NAME}.zip" "${OUTPUT_NAME}${OUTPUT_SUFFIX}"

# Upload the Release
curl "${UPLOAD_URL}?name=${OUTPUT_NAME}.zip" -X POST --data-binary "@${OUTPUT_NAME}.zip" -H 'Content-Type: application/zip' -H "Authorization: Bearer ${GITHUB_TOKEN}"
