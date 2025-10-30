FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o feedback-server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/feedback-server .

EXPOSE 8080

CMD ["./feedback-server"]
