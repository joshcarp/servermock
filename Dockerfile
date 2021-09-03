FROM golang:alpine
WORKDIR /usr/app
ADD . .
RUN go build -o /bin ./cmd/dmt
ENTRYPOINT dmt
