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
git clone --depth 1 "https://github.com/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}" || su -c "!!" || exit 10
cd "${GITHUB_WORKSPACE}" || exit 13

# Host Environment
export GOOS=linux
export GOARCH=amd64

# Go GOPATH!
go get -d . || go get -d . || exit 20

# Convert Icon
go get -u github.com/cratonica/2goarray || su -c "!!" || exit 30
~/go/bin/2goarray Icon main < icon.ico > icon.go

# Embed Icon in Resource
go get -u github.com/akavel/rsrc || su -c "!!" || exit 40
~/go/bin/rsrc -ico icon.ico

# Go Build!
GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc go build -ldflags '-linkmode=internal -H=windowsgui' . && mv psftp.exe /psftp32.exe || exit 50
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags '-linkmode=internal -H=windowsgui' . && mv psftp.exe /psftp64.exe || exit 51
#GOOS=linux GOARCH=386 CC="gcc -m32 -melf_i386" go build -ldflags '-linkmode=internal' . && mv psftp /psftp32 || exit 60
GOOS=linux GOARCH=amd64 CC="gcc -m64" go build -ldflags '-linkmode=internal' . && mv psftp /psftp64 || exit 61

# Exports
if [ "${GITHUB_ACTIONS}" == 'true' ]; then
	echo "::set-env name=OUTPUT_NAMES::$(ls /psftp*)"
fi

# Debugging?
if [ $DEBUG_BUILD -eq 1 ]; then
	/usr/bin/bash
fi
