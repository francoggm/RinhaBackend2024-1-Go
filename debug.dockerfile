FROM golang:1.21.0

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY ./src .

RUN go mod download

RUN go install github.com/cosmtrek/air@latest

EXPOSE 8080 8080

CMD ["air", "-c", ".air.toml"]