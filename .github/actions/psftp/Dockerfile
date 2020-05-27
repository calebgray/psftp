# Ubuntu
FROM ubuntu:latest AS golang

# Install CGO Prerequisites
ENV DEBIAN_FRONTEND noninteractive
RUN apt update && apt upgrade -y && apt install -y build-essential curl git-all pkg-config libxxf86vm-dev libappindicator3-dev gcc-mingw-w64-x86-64

# Install Go
RUN curl https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz | tar xzf - -C /opt/
ENV PATH /opt/go/bin:/root/go/bin:$PATH

# Fork Builder
FROM golang AS builder

# Run the Build!
WORKDIR /root/
COPY . .
ENTRYPOINT /usr/bin/bash
