#!/usr/bin/bash

# Running Locally?
export DEBUG_BUILD=0
if [ "${GITHUB_ACTIONS}" != 'true' ]; then
	# Setup Environment
	export PROJECT_NAME=psftp
	export GITHUB_WORKSPACE=/

	# Interactive!
	echo '`exit` any time to begin build; `exit 1` to hook end of build.'
	/usr/bin/bash
	DEBUG_BUILD=$?
else
	# Dynamic Project Name
	export PROJECT_NAME=$(basename "${GITHUB_REPOSITORY}")
fi

# Sandbox
cd "${GITHUB_WORKSPACE}" || exit 15

# Host Environment
export GOOS=linux
export GOARCH=amd64

# Go GOPATH!
go get -d . || go get -d . || exit 20

# Go Generate!
go generate . || exit 30

# Go RTFM!
cp ./*.md build

# Go Build!
GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc go build -ldflags '-linkmode=internal -H=windowsgui' . && mv "${PROJECT_NAME}.exe" "build/${PROJECT_NAME}32-${PROJECT_VERSION}.exe" || exit 50
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags '-linkmode=internal -H=windowsgui' . && mv "${PROJECT_NAME}.exe" "build/${PROJECT_NAME}64-${PROJECT_VERSION}.exe" || exit 51
#GOOS=linux GOARCH=386 CC="gcc -m32 -melf_i386" go build -ldflags -linkmode=internal . && mv "${PROJECT_NAME}" "build/${PROJECT_NAME}32-${PROJECT_VERSION}-gtk" || exit 60
GOOS=linux GOARCH=amd64 CC="gcc -m64" go build -ldflags -linkmode=internal . && mv "${PROJECT_NAME}" "build/${PROJECT_NAME}64-${PROJECT_VERSION}-gtk" || exit 61

# Exports
if [ "${GITHUB_ACTIONS}" == 'true' ]; then
	echo "::set-env name=BUILD_DIR::$(pwd)/build"
fi

# Debugging?
if [ $DEBUG_BUILD -eq 1 ]; then
	/usr/bin/bash
fi
