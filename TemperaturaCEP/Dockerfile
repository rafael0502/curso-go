FROM golang:1.24.3 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o temperaturacep

FROM scratch
WORKDIR /app
COPY --from=build /app/temperaturacep .
ENTRYPOINT ["./temperaturacep"]