FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-server .

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/chat-server .

ENV PORT=8080
EXPOSE 8080

CMD ["./chat-server"]