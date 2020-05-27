#!/bin/sh

# Before Building
export CC=i686-w64-mingw32-gcc

# Fetch Go Dependencies
go get -v ./... || true

# Go Build!
go build .
