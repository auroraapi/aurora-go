# Use golang image for building
FROM golang:1.9-alpine

# Set GOPATH and download dependencies
ENV GOPATH /go
ARG id_rsa
RUN apk add -U git curl openssh 

# Install go dependencies
RUN mkdir -p ~/.ssh && \
    (echo "$id_rsa" | tr '_' '\n') > ~/.ssh/id_rsa && \
    echo -e "Host *\n\tStrictHostKeyChecking no\n" > ~/.ssh/config && \
    chmod 400 ~/.ssh/* && \
    curl -L https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 -o /go/bin/dep && \
    chmod +x /go/bin/dep

WORKDIR /go/src/aurora-go
COPY Gopkg* /go/src/aurora-go/
RUN eval $(ssh-agent -s) && \
    ssh-add ~/.ssh/id_rsa && \
    dep ensure -v -vendor-only

# Create workdir, copy files, build go binary
COPY . .
RUN go build