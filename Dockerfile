FROM golang:alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o ./bin/app

CMD ["./bin/app"]
