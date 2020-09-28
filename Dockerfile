FROM golang:alpine3.12 AS builder

WORKDIR /go/src/github.com/magiskboy/todo

RUN apk --update add gcc make g++

ADD . .

RUN make

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/github.com/magiskboy/todo/todo .

ENTRYPOINT ./todo web

CMD /bin/sh
