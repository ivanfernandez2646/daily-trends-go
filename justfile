lint:
    go fmt ./...

run:
    go run cmd/main.go

build:
    @CGO_ENABLED=0 go build cmd/main.go

compose:
    docker-compose up -d

down:
    docker-compose down