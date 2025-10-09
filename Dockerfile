FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/backend .
COPY ./scripts/init_db.sh ./scripts/init_db.sh

RUN chmod +x ./scripts/init_db.sh

CMD ["bash", "./scrits/init_db.st"]