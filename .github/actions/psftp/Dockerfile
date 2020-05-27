# Ubuntu
FROM ubuntu:latest AS golang

# Install CGO Prerequisites
ENV DEBIAN_FRONTEND noninteractive
RUN apt update && apt upgrade -y && apt install -y build-essential curl git-all pkg-config libxxf86vm-dev libappindicator3-dev gcc-mingw-w64-x86-64
RUN ln -fs /usr/share/zoneinfo/America/Pacific /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

# Install Go
RUN curl -o /opt/go1.14.3.linux-amd64.tar.gz https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz
RUN tar xf /opt/go1.14.3.linux-amd64.tar.gz -C /opt/
RUN rm /opt/go1.14.3.linux-amd64.tar.gz
ENV PATH /opt/go/bin:$PATH

# New Builder
FROM golang AS builder

# Setup Environment
ENV CC i686-w64-mingw32-gcc

# Fetch Go Dependencies
RUN mkdir -p /root/go/src/github.com/calebgray/
RUN ln -s "$GITHUB_WORKSPACE" /root/go/src/github.com/calebgray/psftp
WORKDIR /root/go/src/github.com/calebgray/psftp
COPY . .
RUN go get -v ./... || true
RUN mkdir /root/go/src/github.com/ffred && cd /root/go/src/github.com/ffred && git clone https://github.com/ffred/guitocons
RUN mkdir -p /root/go/src/golang.org/x && cd /root/go/src/golang.org/x && git clone https://github.com/golang/sys
RUN GOOS=windows GOARCH=386 go build .
