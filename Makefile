build:
	go build -o bin/code-pulse ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

docker-build:
	docker build -t code-pulse .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	rm -rf bin/

.PHONY: build run test docker-build docker-up docker-down clean