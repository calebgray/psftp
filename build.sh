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
	export GITHUB_WORKSPACE=/root
	export PROJECT_ROOT="${HOME}/go/src/build.local/build/${PROJECT_NAME}"

	# Interactive!?
	echo '`exit` any time to begin build; `exit 1` to hook end of build.'
	/usr/bin/bash
	DEBUG_BUILD=$?
else
	# Dynamic Environment
	export PROJECT_VERSION=$(date +%Y.%m.%d).$(echo "${GITHUB_SHA}" | cut -c1-4)
	export PROJECT_NAME=$(basename "${GITHUB_REPOSITORY}")
fi

# Sandbox
mkdir -p "$(dirname "${PROJECT_ROOT}")"
ln -s ${GITHUB_WORKSPACE} "${PROJECT_ROOT}"
cd "${PROJECT_ROOT}" || exit 10

# Go GOPATH!
/usr/bin/bash
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
#GOOS=linux GOARCH=386 go build -ldflags -linkmode=internal . && mv "${PROJECT_NAME}" "build/${PROJECT_NAME}32-${PROJECT_VERSION}-gtk" || exit 60
GOOS=linux GOARCH=amd64 go build -ldflags -linkmode=internal . && mv "${PROJECT_NAME}" "build/${PROJECT_NAME}64-${PROJECT_VERSION}-gtk" || exit 61

# Upload!?
if [ "${GITHUB_ACTIONS}" == 'true' ]; then
	# Establish "Security"
	mkdir -p /root/.ssh
	echo "${UPLOAD_KEY}" > /root/.ssh/id_rsa
	chmod 600 /root/.ssh/id_rsa

	# Establish Even More "Security"
	UPLOAD_HOST=$(echo "${UPLOAD_GIT}" | grep -o '@[^:]*')
	UPLOAD_HOST=${UPLOAD_HOST:1}
	ssh-keyscan -H -t rsa "${UPLOAD_HOST}" > /root/.ssh/known_hosts

	# Identify! Identify!
	git config --global user.email "${UPLOADER_EMAIL}"
	git config --global user.name "${UPLOADER_NAME}"

	# Finally, Upload the Release!
	git clone "${UPLOAD_GIT}" upload || exit 70
	cd upload || exit 71
	git rm -fr build
	mv ../build .
	git add build
	git commit -m "${PROJECT_VERSION}"
	git tag -a "${PROJECT_VERSION}" -m "${GIT_LOG}"
	git push
fi

# Debugging?
if [ $DEBUG_BUILD -eq 1 ]; then
	/usr/bin/bash
fi
