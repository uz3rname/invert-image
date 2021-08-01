TARGET=./bin/app

all:
	$(MAKE) swagger
	$(MAKE) build

build:
	go build -o $(TARGET)

swagger:
	swag init -g ./api/api.go -o ./api/docs

test: build
	go test -v ./api
