#!/bin/sh

# Setup Environment
mkdir -p "/root/go/src/github.com/${GITHUB_ACTOR}"
ln -s "${GITHUB_WORKSPACE}" "/root/go/src/github.com/${GITHUB_REPOSITORY}"
cd "/root/go/src/github.com/${GITHUB_REPOSITORY}" || exit 1

# Before the Build
export CC=i686-w64-mingw32-gcc

# Fetch Go Dependencies
go get -v ./...

# Go Build!
mkdir /root/go/src/github.com/ffred && cd /root/go/src/github.com/ffred && git clone https://github.com/ffred/guitocons
mkdir -p /root/go/src/golang.org/x && cd /root/go/src/golang.org/x && git clone https://github.com/golang/sys
GOOS=windows GOARCH=386 go build .
