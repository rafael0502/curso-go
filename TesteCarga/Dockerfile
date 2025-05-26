# Dockerfile
FROM golang:1.24.3 as builder

WORKDIR /app
COPY . .

RUN go build -o testecarga main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/testecarga .

ENTRYPOINT ["./testecarga"]