#!/usr/bin/bash

# Before Building
export GOOS=linux
export GOARCH=amd64

# Convert Icon
for i in 1 2 3 4 5; do go get -u github.com/cratonica/2goarray && break || sleep 1; done
/root/go/bin/2goarray Icon main < icon.ico > icon.go

# Embed Icon in Resource
for i in 1 2 3 4 5; do go get -u github.com/akavel/rsrc && break || sleep 1; done
/root/go/bin/rsrc -ico icon.ico

# Before Building
export GOOS=windows
export GOARCH=386

# Go Get Deps!
go get -d . || true

# Go Build!
export CC=i686-w64-mingw32-gcc
go build -ldflags "-linkmode=internal -H=windowsgui" .
