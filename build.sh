#!/usr/bin/bash

# Running Locally?
export DEBUG_BUILD=0
if [ "${GITHUB_ACTIONS}" != 'true' ]; then
	# Setup Environment
	export GITHUB_REPOSITORY=calebgray/psftp
	export PROJECT_NAME=$(basename "${GITHUB_REPOSITORY}")
	export GITHUB_WORKSPACE=/${PROJECT_NAME}

	# Interactive!
	echo '`exit` any time to begin build; `exit 1` to hook end of build.'
	/usr/bin/bash
	DEBUG_BUILD=$?
else
	export PROJECT_NAME=$(basename "${GITHUB_REPOSITORY}")
fi

# Clone Sources
git clone --depth 1 "https://github.com/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}" || !! || exit 10
cd "${GITHUB_WORKSPACE}" || exit 13

# Build Tools
export GOOS=linux
export GOARCH=amd64

# Go GOPATH!
go get -d . || go get -d . || exit 20

# Convert Icon
go get -u github.com/cratonica/2goarray || !! || exit 30
~/go/bin/2goarray Icon main < icon.ico > icon.go

# Embed Icon in Resource
go get -u github.com/akavel/rsrc || !! || exit 40
~/go/bin/rsrc -ico icon.ico

# Go Build!
export GOOS=windows
export GOARCH=386
export CC=i686-w64-mingw32-gcc
go build -ldflags "-linkmode=internal -H=windowsgui" .
mv psftp.exe /psftp32.exe

export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
go build -ldflags "-linkmode=internal -H=windowsgui" .
mv psftp.exe /psftp64.exe

export GOOS=linux
export GOARCH=386
export CC=gcc
go build .
mv psftp /psftp32

export GOARCH=amd64
go build .
mv psftp /psftp64

ls /psftp*

# Debugging?
if [ $DEBUG_BUILD -eq 1 ]; then
	/usr/bin/bash
fi
