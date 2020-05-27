#!/usr/bin/env bash

# Before Building
export GOOS=windows
export GOARCH=386
export CC=i686-w64-mingw32-gcc

# Go Get Deps!
rm -fr ~/go && go get -d .
while [ $? -ne 0 ]; do !!; done

# Go Build!
go build .
