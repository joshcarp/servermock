FROM golang:alpine
WORKDIR /usr/app
ADD . .
RUN go build ./cmd/mirror
RUN apk update && apk add ca-certificates && update-ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ./cmd/mirror