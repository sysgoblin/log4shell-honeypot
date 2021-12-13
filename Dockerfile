FROM golang:alpine

WORKDIR /app

COPY extractor extractor
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go mod download

RUN go build -o l4sh .