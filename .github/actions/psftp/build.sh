#!/bin/sh

# Setup Environment
mkdir -p "/github/home/go/src/github.com/${GITHUB_ACTOR}"
ln -s "${GITHUB_WORKSPACE}" "/github/home/go/src/github.com/${GITHUB_REPOSITORY}"
cd "/github/home/go/src/github.com/${GITHUB_REPOSITORY}" || exit 1

# Before Building
export CC=i686-w64-mingw32-gcc

# Fetch Go Dependencies
go get -v ./... || true
go get -u golang.org/x/sys
mkdir /github/home/go/src/github.com/ffred && cd /github/home/go/src/github.com/ffred && git clone https://github.com/ffred/guitocons

# Go Build!
GOOS=windows GOARCH=386 go build .
