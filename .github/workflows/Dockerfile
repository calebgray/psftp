FROM golang:1.14-alpine

RUN apk add --no-cache curl jq git build-base

ADD build.sh /build.sh
CMD /build.sh