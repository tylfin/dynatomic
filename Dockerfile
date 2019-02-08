FROM golang:1.11.4
ENV GODEBUG netdns=cgo

ADD . /go/src/github.com/tylfin/dynatomic
WORKDIR /go/src/github.com/tylfin/dynatomic
