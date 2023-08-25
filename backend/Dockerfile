FROM golang:1.20-alpine as builder

WORKDIR /app

RUN apk add --update gcc musl-dev

COPY . /app
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/tundra-dns ./cmd/main.go

FROM alpine:latest as production
COPY --from=builder /app/tundra-dns .
CMD ["./tundra-dns"]
