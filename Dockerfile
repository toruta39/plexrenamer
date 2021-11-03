FROM golang:1.17.2-alpine3.14 AS builder

WORKDIR /app

COPY . .

RUN go build -o plexrenamer cmd/main.go

FROM alpine:3.14

WORKDIR /app

ARG USER=nonroot
RUN apk add --update sudo
RUN adduser -D $USER \
  && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
  && chmod 0440 /etc/sudoers.d/$USER

COPY --from=builder /app/plexrenamer /app/plexrenamer

VOLUME ["/media"]

USER nonroot

ENTRYPOINT ["/app/plexrenamer", "-from", "/media", "-to", "/media"]