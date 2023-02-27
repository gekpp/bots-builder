FROM golang:1.19-buster AS builder

WORKDIR /go/src/app

COPY . .

RUN go build -o /go/bin/app ./cmd/*.go

WORKDIR /dist
RUN mkdir opt && cp /go/bin/app opt/app
RUN ldd /go/bin/app | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'

FROM alpine:latest

ARG QNR_ID
ARG TOKEN
ARG DATABASE_HOST
ARG DATABASE_PORT=5432
ARG DATABASE_NAME=bots_builder
ARG DATABASE_USER=postgres
ARG DATABASE_PASS=adming
ARG DATABASE_CONNECT_TIMEOUT=15

COPY --from=builder /dist /
WORKDIR /opt

ENTRYPOINT ["/opt/app"]
