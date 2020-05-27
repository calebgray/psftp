# Ubuntu
FROM ubuntu:latest AS golang

# Install CGO Prerequisites
ENV DEBIAN_FRONTEND noninteractive
RUN apt update && apt upgrade -y && apt install -y build-essential curl git-all pkg-config libxxf86vm-dev libappindicator3-dev gcc-mingw-w64-x86-64
RUN dpkg-reconfigure --frontend noninteractive tzdata

# Install Go
RUN curl -o /opt/go1.14.3.linux-amd64.tar.gz https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz
RUN tar xf /opt/go1.14.3.linux-amd64.tar.gz -C /opt/
RUN rm /opt/go1.14.3.linux-amd64.tar.gz
ENV PATH /opt/go/bin:$PATH

# Fork Builder
FROM golang AS builder

# Run the Build!
CMD .github/actions/psftp/build.sh
