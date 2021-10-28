FROM golang:alpine
WORKDIR /usr/app
ADD . .
RUN go build -o /bin ./cmd/servermock
ENTRYPOINT servermock
