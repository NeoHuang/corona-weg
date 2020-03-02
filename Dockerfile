FROM golang:1.13-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o corona-weg .

EXPOSE 8404

CMD ["./corona-weg"]
