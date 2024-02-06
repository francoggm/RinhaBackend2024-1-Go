FROM golang:1.21.0

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/service

EXPOSE 8080 8080

ENTRYPOINT ["build/service"]