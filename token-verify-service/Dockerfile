FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o service-2 ./commander


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/service-2 .

EXPOSE 8000

# Command to run the executable
CMD ["./service-2"]