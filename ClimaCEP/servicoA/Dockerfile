FROM golang:1.24.3

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /servicoa

EXPOSE 8082

CMD ["/servicoa"]