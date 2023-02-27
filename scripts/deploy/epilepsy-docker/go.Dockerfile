FROM golang:1.19-buster AS builder

WORKDIR /go/src/app

COPY . .

RUN go build -o /go/bin/app ./cmd/*.go

FROM alpine:latest

WORKDIR /opt/

ENV APP_PATH="/opt/app"

COPY --from=builder /go/bin/app $APP_PATH

ENTRYPOINT ["/opt/app"]
