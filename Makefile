.PHONY: run seed build test infra-up infra-down docker-up docker-down clean

run:
	go run cmd/api/main.go

seed:
	go run cmd/api/main.go seed

build:
	go build -o bin/super-indo-api cmd/api/main.go

test:
	go test ./... -v -cover

infra-up:
	docker-compose up -d postgres redis

infra-down:
	docker-compose down

docker-up:
	docker-compose --profile fullstack up -d --build

docker-down:
	docker-compose --profile fullstack down

clean:
	rm -rf bin/
