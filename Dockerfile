FROM golang:alpine

WORKDIR /app
COPY . .

RUN go get
RUN go build

CMD ["./invert-image"]
