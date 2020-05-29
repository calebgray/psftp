#!/usr/bin/bash

# Prepare Linux
if [ "${GITHUB_ACTIONS}" == 'true' ]; then
	# Install CGO Prerequisites
	export DEBIAN_FRONTEND=noninteractive
	apt update && apt upgrade -y && apt install -y build-essential curl git-all pkg-config libxxf86vm-dev libappindicator3-dev gcc-mingw-w64-x86-64

	# Install Go
	curl https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz | tar xzf - -C /opt/
	export PATH=/opt/go/bin:~/go/bin:$PATH
fi

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
	mkdir -p /root/.ssh
	echo "${UPLOAD_KEY}" > /root/.ssh/id_rsa
	chmod 600 /root/.ssh/id_rsa

	UPLOAD_HOST=$(echo "${UPLOAD_GIT}" | grep -o '@[^:]*')
	UPLOAD_HOST=${UPLOAD_HOST:1}
	ssh-keyscan -H -t rsa "${UPLOAD_HOST}" > /root/.ssh/known_hosts

	git config --global user.email "${UPLOADER_EMAIL}"
  git config --global user.name "${UPLOADER_NAME}"

	git clone "${UPLOAD_GIT}" upload || exit 70
	cd upload || exit 71
	git rm -fr build
	mv ../build .
	git commit -am "${PROJECT_VERSION}" build
	git push
fi
