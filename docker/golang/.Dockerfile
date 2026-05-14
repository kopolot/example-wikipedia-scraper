FROM golang:1.25 AS builder

WORKDIR /app

COPY app/ .

RUN go mod tidy
RUN make build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin .

RUN apk add chromium

CMD ["./bin/api.bin"]