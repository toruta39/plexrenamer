FROM golang:1.17.2-alpine3.14 AS builder

WORKDIR /app

COPY . .

RUN go build -o plexrenamer cmd/main.go

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/plexrenamer /app/plexrenamer

ENTRYPOINT ["/app/plexrenamer"]