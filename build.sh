#!/usr/bin/bash

# Host Environment
export GOOS=linux
export GOARCH=amd64
export GIT_LOG=$(git log --pretty=oneline)

# Running Locally?
export DEBUG_BUILD=0
if [ "${GITHUB_ACTIONS}" != 'true' ]; then
	# Development Environment
	export PROJECT_VERSION=dev
	export PROJECT_NAME=psftp
	export GITHUB_WORKSPACE=/

	# Interactive!
	echo '`exit` any time to begin build; `exit 1` to hook end of build.'
	/usr/bin/bash
	DEBUG_BUILD=$?
else
	# Dynamic Environment
	export PROJECT_VERSION=$(date +%Y.%m.%d).$(echo "${GITHUB_SHA}" | cut -c1-4)
	export PROJECT_NAME=$(basename "${GITHUB_REPOSITORY}")
fi

# Sandbox
cd "${GITHUB_WORKSPACE}" || exit 15

# Go GOPATH!
go get -d . || go get -d . || exit 20

# Go Generate!
go generate . || exit 30

# All Systems Go...!
mkdir build

# Go RTFM!
cp ./*.md build/

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

# Upload!?
if [ "${GITHUB_ACTIONS}" == 'true' ]; then
	mkdir ~/.ssh
	echo "${UPLOAD_KEY}" > ~/.ssh/id_rsa
	chmod 600 ~/.ssh/id_rsa

	UPLOAD_HOST=$(echo "${UPLOAD_GIT}" | grep -o '@[^:]*')
	UPLOAD_HOST=${UPLOAD_HOST:1}
	ssh-keyscan -H -t rsa "${UPLOAD_HOST}" > ~/.ssh/known_hosts

	ssh -vv git@github.com

	git clone "${UPLOAD_GIT}" upload || exit 70
	cd upload || exit 71
	git rm -fr build
	mv ../build .
	git commit -am "${PROJECT_VERSION}"
	git push
fi
