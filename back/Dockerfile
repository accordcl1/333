FROM golang:1.26.8-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go main_test.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server .

FROM alpine:3.20 AS runtime

WORKDIR /app

COPY --from=builder /server /app/server

EXPOSE 5000

CMD ["/app/server"]
