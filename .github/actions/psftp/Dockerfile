# Ubuntu
FROM ubuntu

# Install CGO Prerequisites
ENV DEBIAN_FRONTEND noninteractive
RUN apt update && apt upgrade -y && apt install -y build-essential curl git-all pkg-config libxxf86vm-dev libappindicator3-dev gcc-mingw-w64-x86-64

# Install Go
RUN curl https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz | tar xzf - -C /opt/
ENV PATH /opt/go/bin:~/go/bin:$PATH

# Run the Build!
COPY ./build.sh .
CMD ./build.sh
