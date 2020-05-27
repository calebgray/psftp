# Ubuntu
FROM ubuntu:latest AS builder

# Install Development Prerequisites
ENV DEBIAN_FRONTEND noninteractive
RUN apt update && apt upgrade -y && apt install -y build-essential curl git-all pkg-config libxxf86vm-dev libappindicator3-dev gcc-mingw-w64-x86-64
RUN ln -fs /usr/share/zoneinfo/America/Pacific /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

# Install Go
RUN wget -P /opt/ https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz
RUN tar xf /opt/go1.14.3.linux-amd64.tar.gz -C /opt/
RUN rm /opt/go1.14.3.linux-amd64.tar.gz
ENV PATH /opt/go/bin:$PATH

# Build Sources
RUN mkdir -p /go/src/github.com/calebgray/psftp
WORKDIR /go/src/github.com/calebgray/psftp
COPY . .
RUN go get -v ./...
RUN GOOS=windows GOARCH=386 go build .
