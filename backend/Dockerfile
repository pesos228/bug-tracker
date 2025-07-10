FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /app/server ./cmd/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .


EXPOSE 8081

CMD ["./server"]