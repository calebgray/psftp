#!/bin/sh

# Setup Environment
mkdir -p "/github/home/go/src/github.com/${GITHUB_ACTOR}"
ln -s "${GITHUB_WORKSPACE}" "/github/home/go/src/github.com/${GITHUB_REPOSITORY}"
cd "/github/home/go/src/github.com/${GITHUB_REPOSITORY}" || exit 1

# Before the Build
export CC=i686-w64-mingw32-gcc

# Fetch Go Dependencies
go get -v ./... || true

# Go Build!
mkdir /github/home/go/src/github.com/ffred && cd /github/home/go/src/github.com/ffred && git clone https://github.com/ffred/guitocons
mkdir -p /github/home/go/src/golang.org/x && cd /github/home/go/src/golang.org/x && git clone https://github.com/golang/sys
GOOS=windows GOARCH=386 go build .
