# ---- Build stage ----
FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ---- Dev stage (with Air for hot reload) ----
FROM golang:1.22
WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY --from=builder /app /app

EXPOSE 8080

CMD ["air", "cmd/main.go"]
