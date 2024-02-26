FROM golang:1.21.0-alpine AS builder

WORKDIR /src

COPY ./src .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o service

FROM scratch

WORKDIR /dist

COPY --from=builder /src/service ./

CMD ["./service"]